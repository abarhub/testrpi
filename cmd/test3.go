package main

import (
	"time"
)
import (
	"machine"
	"tinygo.org/x/drivers/tm1637"
)

func main() {

	tmp2 := machine.I2C0

	println("test", tmp2)
	//println("test", tmp2, machine.GPIO2, machine.GPIO3)

	//tm := tm1637.New(machine.D2, machine.D3, 7) // clk, dio, brightness
	tm := tm1637.New(2, 3, 7) // clk, dio, brightness
	tm.Configure()

	tm.ClearDisplay()

	tm.DisplayText([]byte("Tiny"))
	time.Sleep(time.Millisecond * 1000)

	tm.ClearDisplay()

	tm.DisplayChr(byte('G'), 1)
	tm.DisplayDigit(0, 2) // looks like O
	time.Sleep(time.Millisecond * 1000)

	tm.DisplayClock(12, 59, true)

	for i := uint8(0); i < 8; i++ {
		tm.Brightness(i)
		time.Sleep(time.Millisecond * 200)
	}

	i := int16(0)
	for {
		tm.DisplayNumber(i)
		i++
		time.Sleep(time.Millisecond * 50)
	}

}

func main0() {
	println("demarrage")

	machine.I2C0.Configure(machine.I2CConfig{})
	//sensor := bmp180.New(machine.I2C0)
	//sensor.Configure()
	//
	//connected := sensor.Connected()
	//if !connected {
	//	println("BMP180 not detected")
	//	return
	//}
	println("BMP180 detected")

	tmp := tm1637.New(5, 4, 7)

	println("affichage")

	tmp.DisplayNumber(1234)

	println("attente")
	time.Sleep(2 * time.Second)

	println("fin")

	//machine.I2C0.Configure(machine.I2CConfig{})
	//sensor := bmp180.New(machine.I2C0)
	//sensor.Configure()
	//
	//connected := sensor.Connected()
	//if !connected {
	//	println("BMP180 not detected")
	//	return
	//}
	//println("BMP180 detected")
	//
	//for {
	//	temp, _ := sensor.ReadTemperature()
	//	println("Temperature:", float32(temp)/1000, "°C")
	//
	//	pressure, _ := sensor.ReadPressure()
	//	println("Pressure", float32(pressure)/100000, "hPa")
	//
	//	time.Sleep(2 * time.Second)
	//}
}
