package main

import (
	// "fmt"
	"log"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
	"tinygo.org/x/tinyfont"

	// "tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/gophers"
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

	time.Sleep(2 * time.Second)
	log.Println("Starting...")

	log.Println("Configure SPI...")
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
		SCK:       clk,
		SDO:       din,
		// SDI:       sdi,
	})

	log.Println("new epd4in2")
	display = epd4in2.New(machine.SPI0, cs, dc, rst, busy)
	log.Println("Configure epd4in2")
	display.Configure(epd4in2.Config{})

	// clearDisplay(&display)
	// displayCheckerBoard(&display)

	log.Println("clearDisplay()")
	clearDisplay(&display)
	log.Println("fontExamples()")
	fontExamples(&display)

	log.Println("Waiting for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Println("You could remove power now")
}

func showRect(x int16, y int16, w int16, h int16, c color.RGBA) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			display.SetPixel(i, j, c)
		}
	}
}

func clearDisplay(display *epd4in2.Device) {

	log.Println("ClearBuffer")
	display.ClearBuffer()
	log.Println("ClearDisplay")
	display.ClearDisplay()
	log.Println("WaitUntilIdle")
	display.WaitUntilIdle()
	log.Println("Waiting for 3 seconds")
	time.Sleep(3 * time.Second)
}

func displayCheckerBoard(display *epd4in2.Device) {
	black := color.RGBA{1, 1, 1, 255}

	log.Println("Prep checkered board")
	// Show a checkered board
	for i := int16(0); i < 49; i++ {
		for j := int16(0); j < 29; j++ {
			if (i+j)%2 == 0 {
				showRect(i*8, j*10, 8, 10, black)
			}
		}
	}

	log.Println("Waiting for 3 seconds")
	time.Sleep(3 * time.Second)

	log.Println("Display()")
	display.Display()
	log.Println("WaitUntilIdle()")
	display.WaitUntilIdle()

	log.Println("Waiting for 3 seconds")
	time.Sleep(3 * time.Second)
}

func fontExamples(display *epd4in2.Device) {

	black := color.RGBA{1, 1, 1, 255}
	// white := color.RGBA{0, 0, 0, 255}

	//tinyfont.WriteLineRotated(display, &freemono.Bold9pt7b, 15, 20, "a", black, tinyfont.NO_ROTATION)

	time.Sleep(3 * time.Second)

	// showRect(0, 22, 52, 20, black)
	// showRect(52, 22, 52, 20, white)

	// display.Display()
	// display.WaitUntilIdle()

	// log.Println("Waiting for 5 seconds")
	// time.Sleep(5 * time.Second)

	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 85, 26, "World!", white, tinyfont.ROTATION_180)
	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 55, 60, "@tinyGolang", black, tinyfont.ROTATION_90)

	// tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50, "ABCDEFG\nHIJKLMN\nOPQRSTU", black, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50,  "ABCDEFG\nHIJKLMN\nOPQRSTU\nH", black, tinyfont.NO_ROTATION)

	// tinyfont.WriteLineColorsRotated(&display, &freemono.Bold9pt7b, 45, 180, "tinyfont", []color.RGBA{white, black}, tinyfont.ROTATION_270)

	log.Println("Waiting for 3 seconds")

	log.Println("Display()")
	display.Display()
	log.Println("WaitUntilIdle()")
	display.WaitUntilIdle()

	log.Println("Waiting for 3 seconds")
	time.Sleep(3 * time.Second)
}
