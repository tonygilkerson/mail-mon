package main

import (
	"fmt"
	"log"
	"machine"
	"time"

	"image/color"

	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/drivers/tone"
)

func main() {
	//
	// Named PINs
	//
	var dspKey2 machine.Pin = machine.GP2
	var dspKey3 machine.Pin = machine.GP3
	var buzzerPin machine.Pin = machine.GP7 
	var dspDC machine.Pin = machine.GP8
	var dspCS machine.Pin = machine.GP9
	var dspSCK machine.Pin = machine.GP10
	var dspSDO machine.Pin = machine.GP11
	var dspReset machine.Pin = machine.GP12
	var dspLite machine.Pin = machine.GP13
	var dspKey0 machine.Pin = machine.GP15
	var dspKey1 machine.Pin = machine.GP17
	var dspSDI machine.Pin = machine.GP28

	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsp.RunLight(led, 10)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// PWM for tone alarm
	//
	buzzer, err := tone.New(machine.PWM3, buzzerPin)
	if err != nil {
		log.Panicln("failed to configure PWM")
	}
	soundSiren(buzzer)

	//
	// Display Buttons
	//
	chKeyPress := make(chan string, 1)
	dspKey0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dspKey1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dspKey2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dspKey3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	dspKey0.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		fmt.Println("key0 - RADriverCmd")
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chKeyPress <- "key0":
		default:
		}
	})

	dspKey1.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		fmt.Println("key1 interrupt")
		select {
		case chKeyPress <- "key1":
		default:
		}
	})

	dspKey2.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		fmt.Println("key2 interrupt")
		select {
		case chKeyPress <- "key2":
		default:
		}
	})

	dspKey3.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		fmt.Println("key3 interrupt")
		select {
		case chKeyPress <- "key3":
		default:
		}
	})

	//
	// Display
	//
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8000000,
		LSBFirst:  false,
		Mode:      0,
		SCK:       dspSCK,
		SDO:       dspSDO,
		SDI:       dspSDI, // I don't think this is actually used for LCD, just assign to any open pin
	})

	display := st7789.New(machine.SPI1,
		dspReset, // TFT_RESET
		dspDC,    // TFT_DC
		dspCS,    // TFT_CS
		dspLite)  // TFT_LITE

	display.Configure(st7789.Config{
		// With the display in portrait and the usb socket on the left and in the back
		// the actual width and height are switched width=320 and height=240
		Width:        240,
		Height:       320,
		Rotation:     st7789.ROTATION_90,
		RowOffset:    0,
		ColumnOffset: 0,
		FrameRate:    st7789.FRAMERATE_111,
		VSyncLines:   st7789.MAX_VSYNC_SCANLINES,
	})

	//
	// Start
	//
	log.Printf("start")

	width, height := display.Size()
	log.Printf("width: %v, height: %v\n", width, height)

	// red := color.RGBA{126, 0, 0, 255} // dim
	red := color.RGBA{255, 0, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	// green := color.RGBA{0, 255, 0, 255}

	lastTakenMedsAt := time.Now()
	screenOnAt := time.Now()

	for {

		select {
		case key := <-chKeyPress:

			log.Println("key channel message %s", key)
			
			switch key {

			case "key0":
				log.Println("Key0 - Add 1hr")
				lastTakenMedsAt = lastTakenMedsAt.Add(time.Hour)
			case "key1":
				log.Println("Key1 - Add 30min")
				lastTakenMedsAt = lastTakenMedsAt.Add(time.Minute * 30)
			case "key2":
				log.Println("Key2 - Do reset/took meds")
				lastTakenMedsAt = time.Now()
			case "key3":
				log.Println("Key3 pressed - Subtract 1hr")
				lastTakenMedsAt = lastTakenMedsAt.Add(time.Hour * -1)
			}

		default:
			time.Sleep(time.Millisecond * 50)
			log.Println(".")
		}

		lastTakenMedsDuration := time.Since(lastTakenMedsAt)
		age := fmt.Sprintf("%1.2fh", lastTakenMedsDuration.Hours())
		ageString := fmt.Sprintf("Last taken:\n%s hours ago", age)


		cls(&display)
		// tinyfont.WriteLine(&display,&freemono.Regular12pt7b,10,20,"123456789-123456789-x",red)
		tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 10, 20, ageString, red)
		time.Sleep(time.Second * 3)

		//test
		// soundSiren(buzzer)

		screenOnDuration := time.Since(screenOnAt)
		if screenOnDuration.Minutes() > 1 {
			//turn off the screen
		}

	}

}

func paintScreen(c color.RGBA, d *st7789.Device, s int16) {
	var x, y int16
	for y = 0; y < 240; y = y + s {
		for x = 0; x < 320; x = x + s {
			d.FillRectangle(x, y, s, s, c)
		}
	}
}

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
	fmt.Printf("FillScreen(black)\n")
}

func soundSiren(buzzer tone.Speaker) {
	for i := 0; i < 5; i++ {
		log.Println("nee")
		buzzer.SetNote(tone.B5)
		time.Sleep(time.Second / 2)

		log.Println("naw")
		buzzer.SetNote(tone.A5)
		time.Sleep(time.Second / 2)

	}
	buzzer.Stop()
}
