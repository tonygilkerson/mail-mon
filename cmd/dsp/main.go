package main

import (
	"fmt"
	"log"
	"machine"
	"runtime"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
	"tinygo.org/x/drivers/sx127x"
)

const (
	HEARTBEAT_DURATION_SECONDS = 300
	TXRX_LOOP_TICKER_DURATION_SECONDS = 9
)


/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	//
	// Named PINs
	//

	var enLora machine.Pin = machine.GP15
	var sdiLora machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var csLora machine.Pin = machine.GP17
	var sckLora machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var sdoLora machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN
	var rstLora machine.Pin = machine.GP20
	var dio0Lora machine.Pin = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	var dio1Lora machine.Pin = machine.GP22 // (GP22--G1)I don't now what this does but it seems to need to be connected

	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// 	Setup Lora
	//
	var loraRadio *sx127x.Device
	txQ := make(chan string, 250) // I would hope the channel size would never be larger than ~4 so 250 is large
	rxQ := make(chan string) // this app currently does not do anything with messages received

	radio := road.SetupLora(*machine.SPI0, enLora, rstLora, csLora, dio0Lora, dio1Lora, sckLora, sdoLora, sdiLora, loraRadio, &txQ, &rxQ, 0, 0, TXRX_LOOP_TICKER_DURATION_SECONDS, road.TxOnly)


	// Launch go routines
	go radio.LoraRxTx()

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int
	
	for range ticker.C {
		
		log.Printf("------------------MainLoopHeartbeat-------------------- %v", count)
		count += 1

		//
		// Send Heartbeat to Tx queue
		//
		txQ <- "DisplayMainLoopHeartbeat"


		runtime.Gosched()
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print("run-")
	}

}


