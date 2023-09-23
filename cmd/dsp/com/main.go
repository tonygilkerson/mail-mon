//
// The program is the communication component for the display
//
// It manages the LORA RX/TX cycle.  Messages received via LORA are 
// send to the epaper component via the UART message bus
//
package main

import (
	"log"
	"machine"
	"runtime"
	"time"

	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"github.com/tonygilkerson/mbx-iot/internal/road"
	"github.com/tonygilkerson/mbx-iot/pkg/iot"
	"tinygo.org/x/drivers/sx127x"
)

const (
	// HEARTBEAT_DURATION_SECONDS        = 300
	HEARTBEAT_DURATION_SECONDS        = 11
	TXRX_LOOP_TICKER_DURATION_SECONDS = 9
)

/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	//
	// Named PINs
	//

	var loraEn machine.Pin = machine.GP15
	var loraSdi machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var loraCs machine.Pin = machine.GP17
	var loraSck machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var loraSdo machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN
	var loraRst machine.Pin = machine.GP20
	var loraDio0 machine.Pin = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	var loraDio1 machine.Pin = machine.GP22 // (GP22--G1) I don't now what this does but it seems to need to be connected

	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsp.RunLight(led, 10)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// 	Setup Lora
	//
	var loraRadio *sx127x.Device
	txQ := make(chan string, 250) // I would hope the channel size would never be larger than ~4 so 250 is large
	rxQ := make(chan string, 250)

	log.Println("Setup LORA")
	radio := road.SetupLora(*machine.SPI0, loraEn, loraRst, loraCs, loraDio0, loraDio1, loraSck, loraSdo, loraSdi, loraRadio, &txQ, &rxQ, 0, 10_000, TXRX_LOOP_TICKER_DURATION_SECONDS, road.TxRx)

	//
	// go routines
	//
	// DEVTODO - it works as a go routine but I dont want it to run as a go routine
	go radio.LoraRxTxRunner()

	//
	// Main loop
	//
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int

	for range ticker.C {

		log.Printf("------------------MainLoopHeartbeat-------------------- %v", count)
		count += 1

		//
		// Send Heartbeat to Tx queue
		//
		txQ <- iot.DspMainLoopHeartbeat
		dsp.RunLight(led, 2)

		//
		// Do a Lora Rx Tx cycle
		//
		// radio.LoraRxTx()

		//
		// Consume any messages received
		//
		rxQConsumer(&rxQ)

		//
		// Let someone else have a turn
		//
		runtime.Gosched()
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func rxQConsumer(rxQ *chan string) {
	var msgBatch string

	for len(*rxQ) > 0 {

		msgBatch = <- *rxQ
		log.Printf("Message batch: [%v]", msgBatch)

		messages := road.SplitMessageBatch(msgBatch)
		for _, msg := range messages {
			log.Printf("dsp.rxQConsumer: Message: [%v]", msg)
		}

	}
}
