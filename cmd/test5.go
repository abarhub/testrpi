package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/devices/v3/tm1637"
	"periph.io/x/host/v3"
	"strconv"
	"strings"
	"time"
)

const HIDE = "x"
const AFFICHE_HEURE = "HEURE"
const AFFICHE_RIEN = "RIEN"
const AFFICHE_MINUTEUR = "MINUTEUR"

// const INTENSITE=tm1637.Brightness10
const INTENSITE = tm1637.Brightness4

type Action struct {
	action string
	heure  int
	minute int
}

type Heure struct {
	afficheHeure bool
	fin          bool
	heure        int
	minute       int
}

var messages = make(chan Heure)
var horlogeEnd = false
var action = make(chan Action)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format("02 Jan 2006 15:04:05 MST"))
}

func initAfficheur() *tm1637.Dev {
	log.Print("initialisation de l'affichage ...")
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	clk := gpioreg.ByName("GPIO5")
	data := gpioreg.ByName("GPIO4")
	if clk == nil || data == nil {
		log.Fatal("Failed to find pins")
	}
	dev, err := tm1637.New(clk, data)
	if err != nil {
		log.Fatalf("failed to initialize tm1637: %v", err)
	}
	if err := dev.SetBrightness(INTENSITE); err != nil {
		log.Fatalf("failed to set brightness on tm1637: %v", err)
	}
	log.Print("initialisation de l'affichage ok")
	return dev
}

func affiche() {

	dev := initAfficheur()

	log.Print("démarrage de la boucle d'affichage ...")

	for {
		msg := <-messages
		if msg.fin {
			break
		} else if !msg.afficheHeure {
			if err := dev.Halt(); err != nil {
				log.Fatalf("failed to halt to tm1637: %v", err)
			}
		} else if msg.afficheHeure && msg.heure >= 0 && msg.heure <= 99 && msg.minute >= 0 && msg.minute <= 99 {
			hours := msg.heure
			minutes := msg.minute
			if err := dev.SetBrightness(INTENSITE); err != nil {
				log.Fatalf("failed to write to tm1637: %v", err)
			}
			if _, err := dev.Write(tm1637.Clock(hours, minutes, true)); err != nil {
				log.Fatalf("failed to write to tm1637: %v", err)
			}
		}
	}

}

func horloge() {
	action <- Action{AFFICHE_HEURE, 0, 0}
}

func minuteur(minute, secondes int) {
	action <- Action{AFFICHE_MINUTEUR, minute, secondes}
}

func arret() {
	action <- Action{AFFICHE_RIEN, 0, 0}
}

func boucleEvenement() {

	var actionPrecedante Action
	var minutes, secondes int
	var begin, end time.Time
	var difference time.Duration

	log.Print("démarrage de la boucle d'evenement ...")

	for {

		nouveau := false
		var actionSelectionnee Action

		select {
		case actionCourante := <-action:
			actionSelectionnee = actionCourante
			nouveau = true
			log.Printf("nouvelle action: %v", actionSelectionnee)
		case <-time.After(1000 * time.Millisecond):
			actionSelectionnee = actionPrecedante
		}

		if actionSelectionnee.action == AFFICHE_HEURE {
			now := time.Now()
			messages <- Heure{afficheHeure: true, heure: now.Hour(), minute: now.Minute()}
		} else if actionSelectionnee.action == AFFICHE_MINUTEUR {
			if nouveau {
				minutes = actionSelectionnee.heure
				secondes = actionSelectionnee.minute
				begin = time.Now()
				end = begin.Add(time.Minute*time.Duration(minutes) + time.Second*time.Duration(secondes))
			} else {

			}
			now := time.Now()
			if now.After(end) || now.Equal(end) {
				messages <- Heure{afficheHeure: true, heure: 0, minute: 0}
			} else {
				diff := end.Sub(now)
				if nouveau || int(diff.Seconds()) != int(difference.Seconds()) {
					minutes := int(diff.Minutes())
					secondes := int(diff.Seconds()) - int(diff.Minutes())*60
					messages <- Heure{afficheHeure: true, heure: minutes, minute: secondes}
					difference = diff
				}
			}
		} else if actionSelectionnee.action == AFFICHE_RIEN {
			if nouveau {
				messages <- Heure{afficheHeure: false}
			}
		}

		//time.Sleep(500 * time.Millisecond)

		actionPrecedante = actionSelectionnee
	}
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "horloge") {
		log.Print("horloge")
		//go func() { horloge() }()
		horloge()
	} else if strings.HasSuffix(r.URL.Path, "minuteur") {
		log.Print("minuteur")
		if r.URL.Query().Has("time") {
			time := r.URL.Query().Get("time")
			log.Print("time", time)
			tab := strings.Split(time, ":")
			log.Print("tab:", tab)
			if len(tab) == 3 {
				minutes, err := strconv.Atoi(tab[1])
				if err == nil {
					secondes, err := strconv.Atoi(tab[2])
					if err == nil && minutes >= 0 && secondes >= 0 && !(minutes == 0 && secondes == 0) {
						minuteur(minutes, secondes)
					} else {
						log.Print("erreur", err)
					}
				} else {
					log.Print("erreur", err)
				}
			}
		}
	} else if strings.HasSuffix(r.URL.Path, "arret") {
		log.Print("arret")
		arret()
	}
	fmt.Fprint(w, "Action")
}

func main() {
	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/static/", fs)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/api/action/horloge", actionHandler)
	http.HandleFunc("/api/action/minuteur", actionHandler)
	http.HandleFunc("/api/action/arret", actionHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	go func() { affiche() }()

	go func() { boucleEvenement() }()

	horloge()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Print("Fin signal")
			arret()
			time.Sleep(10 * time.Second)
			os.Exit(0)
		}
	}()

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Fin")
}
