package main

import (
	"log"
	"machine"
	"runtime"
	"time"

	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"github.com/tonygilkerson/mbx-iot/internal/umsg"
	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
)

const (
	SENDER_ID = "dsp.epaper"
	// HEARTBEAT_DURATION_SECONDS        = 300
	HEARTBEAT_DURATION_SECONDS = 11
)

var display epd4in2.Device

func main() {

	//
	// Named PINs
	//
	var uartInTx machine.Pin = machine.GP0  // UART0
	var uartInRx machine.Pin = machine.GP1  // UART0
	var uartOutTx machine.Pin = machine.GP4 // UART1
	var uartOutRx machine.Pin = machine.GP5 // UART1

	var dc machine.Pin = machine.GP11   // pin15
	var rst machine.Pin = machine.GP12  // pin16
	var busy machine.Pin = machine.GP13 // pin17
	var cs machine.Pin = machine.GP17   // pin22
	var clk machine.Pin = machine.GP18  // pin24 machine.SPI0_SCK_PIN
	var din machine.Pin = machine.GP19  // pin25 machine.SPI0_SDO_PIN

	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// UARTs
	//
	uartIn := machine.UART0
	uartOut := machine.UART1

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsp.RunLight(led, 10)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	/////////////////////////////////////////////////////////////////////////////
	// Broker
	/////////////////////////////////////////////////////////////////////////////

	fooCh := make(chan umsg.FooMsg)
	barCh := make(chan umsg.BarMsg)

	mb := umsg.NewBroker(
		SENDER_ID,

		uartIn,
		uartInTx,
		uartInRx,

		uartOut,
		uartOutTx,
		uartOutRx,

		fooCh,
		barCh,
	)
	log.Printf("[main] - configure message broker\n")
	mb.Configure()

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

	for {
		receiveFooTest(&mb, fooCh)
		runtime.Gosched()
		time.Sleep(1 * time.Second)
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func receiveFooTest(mb *umsg.MsgBroker, fooCh chan umsg.FooMsg) {

	var found bool = false
	var msg umsg.FooMsg

	// Non-blocking ch read that will timeout... boom!
	boom := time.After(1000 * time.Millisecond)
	// for {
		select {
		case msg = <-fooCh:
			found = true
		case <-boom:
			log.Printf("dsp.epaper.receiveFooTest: Boom! timeout waiting for message\n")
			break
		default:
			log.Printf(".")
			runtime.Gosched()
			time.Sleep(50 * time.Millisecond)
		}

		// log.Printf("dsp.epaper.receiveFooTest: found1: %v\n", found)
		// if found {
		// 	break
		// }
		// log.Printf("dsp.epaper.receiveFooTest: found2: %v\n", found)
	// }

	log.Printf("dsp.epaper.receiveFooTest: ******************************************************************\n")
	if found {
		log.Printf("dsp.epaper.receiveFooTest: SUCCESS, msg: [%v]\n", msg)
	} else {
		log.Printf("dsp.epaper.receiveFooTest: FAIL, did not receive message.")
	}
	log.Printf("dsp.epaper.receiveFooTest: ******************************************************************\n")

}
