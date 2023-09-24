package main

import (
	"log"
	"periph.io/x/host/v3"
	//"tinygo.org/x/drivers/tm1637"
)

import (
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/devices/v3/tm1637"
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	//clk := gpioreg.ByName("GPIO6")
	clk := gpioreg.ByName("GPIO5")
	//data := gpioreg.ByName("GPIO12")
	data := gpioreg.ByName("GPIO4")
	if clk == nil || data == nil {
		log.Fatal("Failed to find pins")
	}
	dev, err := tm1637.New(clk, data)
	if err != nil {
		log.Fatalf("failed to initialize tm1637: %v", err)
	}
	if err := dev.SetBrightness(tm1637.Brightness10); err != nil {
		log.Fatalf("failed to set brightness on tm1637: %v", err)
	}
	if _, err := dev.Write(tm1637.Clock(12, 00, true)); err != nil {
		log.Fatalf("failed to write to tm1637: %v", err)
	}
}
