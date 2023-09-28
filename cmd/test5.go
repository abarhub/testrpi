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
	action    string
	heure     int
	minute    int
	intensite int
}

type Heure struct {
	afficheHeure bool
	fin          bool
	heure        int
	minute       int
	intensite    int
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
			intensite := msg.intensite
			intensite2 := INTENSITE
			if intensite > 0 {
				if intensite == 1 {
					intensite2 = tm1637.Brightness1
				} else if intensite == 2 {
					intensite2 = tm1637.Brightness2
				} else if intensite == 3 {
					intensite2 = tm1637.Brightness4
				} else if intensite == 4 {
					intensite2 = tm1637.Brightness10
				} else if intensite == 5 {
					intensite2 = tm1637.Brightness11
				} else if intensite == 6 {
					intensite2 = tm1637.Brightness12
				} else if intensite == 7 {
					intensite2 = tm1637.Brightness13
				} else if intensite == 8 {
					intensite2 = tm1637.Brightness14
				}
			}
			log.Print("intensite:", intensite2, intensite)
			if err := dev.SetBrightness(intensite2); err != nil {
				log.Fatalf("failed to write to tm1637: %v", err)
			}
			//var heure=tm1637.Clock(hours, minutes, true)
			var heure = clock(hours, minutes, true)
			if _, err := dev.Write(heure); err != nil {
				log.Fatalf("failed to write to tm1637: %v", err)
			}
		}
	}

}

// Hex digits from 0 to F.
var digitToSegment = []byte{
	0x3f, 0x06, 0x5b, 0x4f, 0x66, 0x6d, 0x7d, 0x07, 0x7f, 0x6f, 0x77, 0x7c, 0x39, 0x5e, 0x79, 0x71,
}

func clock(hour, minute int, showDots bool) []byte {
	heure := hour / 10
	heure2 := hour % 10
	minute2 := minute / 10
	minute3 := minute % 10
	seg := make([]byte, 4)
	if heure > 0 {
		seg[0] = byte(digitToSegment[heure])
	}
	seg[1] = byte(digitToSegment[heure2])
	seg[2] = byte(digitToSegment[minute2])
	seg[3] = byte(digitToSegment[minute3])
	if showDots {
		seg[1] |= 0x80
	}
	return seg[:]
}

func horloge(intensite int) {
	action <- Action{action: AFFICHE_HEURE, intensite: intensite}
}

func minuteur(minute, secondes int) {
	action <- Action{action: AFFICHE_MINUTEUR, heure: minute, minute: secondes}
}

func arret() {
	action <- Action{action: AFFICHE_RIEN}
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
			messages <- Heure{afficheHeure: true, heure: now.Hour(), minute: now.Minute(), intensite: actionSelectionnee.intensite}
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
		intensite := 0
		if r.URL.Query().Has("intensite") {
			intense, err := strconv.Atoi(r.URL.Query().Get("intensite"))
			if err != nil {
				log.Print("erreur", err)
			} else {
				intensite = intense
			}
		}
		horloge(intensite)
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

	horloge(0)

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
