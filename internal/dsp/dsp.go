package dsp

import (
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
	// "tinygo.org/x/tinyfont/freemono"

)

type Content struct {
	name string
}

func RunLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print("run-")
	}

}


func ClearDisplay(display *epd4in2.Device) {

	display.ClearBuffer()
	display.ClearDisplay()
	display.WaitUntilIdle()
	log.Println("[ClearDisplay] Waiting for 3 seconds")
	time.Sleep(3 * time.Second)

}

func FontExamples(display *epd4in2.Device) {

	black := color.RGBA{1, 1, 1, 255}
	// white := color.RGBA{0, 0, 0, 255}


	time.Sleep(3 * time.Second)


	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 85, 26, "World!", white, tinyfont.ROTATION_180)
	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 55, 60, "@tinyGolang", black, tinyfont.ROTATION_90)

	// tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50, "ABCDEFG\nHIJKLMN\nOPQRSTU", black, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50,  "ABCDEFG\nHIJKLMN\nOPQRSTU\nH", black, tinyfont.NO_ROTATION)

	// tinyfont.WriteLineColorsRotated(&display, &freemono.Bold9pt7b, 45, 180, "tinyfont", []color.RGBA{white, black}, tinyfont.ROTATION_270)


	log.Println("[FontExamples] Display()")
	display.Display()

	log.Println("[FontExamples] WaitUntilIdle()")
	display.WaitUntilIdle()

}