package main

import (
	// "fmt"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
	"github.com/tonygilkerson/mbx-iot/internal/dsp"
)

var display epd4in2.Device

func main() {

	//
	// Named PINs
	//
	var dc machine.Pin = machine.GP11 		// pin15
	var rst machine.Pin = machine.GP12 		// pin16
	var busy machine.Pin = machine.GP13 	// pin17
	var cs machine.Pin = machine.GP17 		// pin22
	var clk machine.Pin = machine.GP18 		// pin24 machine.SPI0_SCK_PIN
	var din machine.Pin = machine.GP19 		// pin25 machine.SPI0_SDO_PIN

	time.Sleep(1 * time.Second)
	log.Println("[main] Starting...")

	//
	// SPI
	//
	log.Println("[main] Configure SPI...")
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
		SCK:       clk,
		SDO:       din,
		// SDI:       sdi,
	})

	//
	// Display
	//
	log.Println("[main] new epd4in2")
	display = epd4in2.New(machine.SPI0, cs, dc, rst, busy)
	display.Configure(epd4in2.Config{})

	log.Println("[main] clearDisplay()")
	dsp.ClearDisplay(&display)

	log.Println("[main] fontExamples()")
	dsp.FontExamples(&display)

	// log.Println("[main] Waiting for 5 seconds")
	// time.Sleep(5 * time.Second)

	log.Println("You could remove power now")
}






