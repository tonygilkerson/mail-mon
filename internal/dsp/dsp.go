package dsp

import (
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
	"tinygo.org/x/tinyfont/freemono"

)

//
// Content for the Display
//
type Content struct {
	isDirty bool
	name string
	gatewayMainLoopHeartbeatStatus string
}

//
// NewContent
//
func NewContent() Content {

	content := Content{
		isDirty: true,
		name: "Mailbox IOT",
		gatewayMainLoopHeartbeatStatus: "initial",
	}

	return content
}


func (content *Content) SetGatewayMainLoopHeartbeatStatus(s string){
	if s != content.gatewayMainLoopHeartbeatStatus {
		content.isDirty = true
		content.gatewayMainLoopHeartbeatStatus = s
	}
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
	print("\n")

}


func ClearDisplay(display *epd4in2.Device) {

	display.ClearBuffer()
	display.ClearDisplay()
	display.WaitUntilIdle()
	log.Println("internal.dsp.ClearDisplay: Waiting for 3 seconds")
	time.Sleep(3 * time.Second)

}

func FontExamples(display *epd4in2.Device) {

	black := color.RGBA{1, 1, 1, 255}
	// white := color.RGBA{0, 0, 0, 255}


	time.Sleep(3 * time.Second)


	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 85, 26, "World!", white, tinyfont.ROTATION_180)
	// tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 55, 60, "@tinyGolang", black, tinyfont.ROTATION_90)

	// tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50, "ABCDEFG\nHIJKLMN\nOPQRSTU", black, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50,  "ABCDEFG\nHIJKLMN\nOPQRSTU\nHH", black, tinyfont.NO_ROTATION)

	// tinyfont.WriteLineColorsRotated(&display, &freemono.Bold9pt7b, 45, 180, "tinyfont", []color.RGBA{white, black}, tinyfont.ROTATION_270)


	log.Println("internal.dsp.FontExamples: Display()")
	display.Display()

	log.Println("internal.dsp.FontExamples: WaitUntilIdle()")
	display.WaitUntilIdle()
	log.Println("internal.dsp.FontExamples: WaitUntilIdle() done.")

}

func  (content *Content) DisplayContent(display *epd4in2.Device) {


	log.Println("internal.dsp.DisplayContent: sleep for a bit!")

	black := color.RGBA{1, 1, 1, 255}
	time.Sleep(3 * time.Second)

	// tinyfont.WriteLineRotated(display, &gophers.Regular58pt, 40, 50,  "HH", black, tinyfont.NO_ROTATION)
	tinyfont.WriteLineRotated(display, &freemono.Bold9pt7b, 30, 50,  "Gateway Heartbeat: "+content.gatewayMainLoopHeartbeatStatus, black, tinyfont.NO_ROTATION)

	log.Println("internal.dsp.DisplayContent: Display()")
	display.Display()

	log.Println("internal.dsp.DisplayContent: WaitUntilIdle()")
	display.WaitUntilIdle()
	log.Println("internal.dsp.DisplayContent: WaitUntilIdle() done.")

}