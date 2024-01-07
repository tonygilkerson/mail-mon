package main

import (
	"log"
	"machine"
	"runtime"
	"time"

	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"github.com/tonygilkerson/mbx-iot/internal/umsg"
	"github.com/tonygilkerson/mbx-iot/pkg/iot"
	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
)

const (
	SENDER_ID = "dsp.epaper"
	HEARTBEAT_DURATION_SECONDS = 300

)

var display epd4in2.Device

func main() {

	//
	// Named PINs
	//
	var uartInTx machine.Pin = machine.GP0 // UART0
	var uartInRx machine.Pin = machine.GP1 // UART0

	var mbxDoorOpenedAckBtn machine.Pin = machine.GP2 // Acknoleges the fact that we got mail, make alerts turn off
	var requestBtn machine.Pin = machine.GP3          // System request, will cycle the heart beat loop and refresh status

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
	// Buttons
	//
	mbxDoorOpenedAckBtnCh := make(chan string, 1)
	mbxDoorOpenedAckBtn.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	mbxDoorOpenedAckBtn.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case mbxDoorOpenedAckBtnCh <- "rise":
		default:
		}

	})

	requestBtnCh := make(chan string, 1)
	requestBtn.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	requestBtn.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case requestBtnCh <- "rise":
		default:
		}

	})

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
	statusCh := make(chan umsg.StatusMsg, 5)

	mb := umsg.NewBroker(
		SENDER_ID,

		uartIn,
		uartInTx,
		uartInRx,

		uartOut,
		uartOutTx,
		uartOutRx,

		fooCh,
		statusCh,
	)
	log.Printf("dsp.com.main: configure message broker\n")
	mb.Configure()

	//
	// SPI
	//
	log.Println("dsp.com.main: Configure SPI...")
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
	log.Println("dsp.com.main: new epd4in2")
	display = epd4in2.New(machine.SPI0, cs, dc, rst, busy)
	display.Configure(epd4in2.Config{})
	content := dsp.NewContent()

	//
	//  Main loop
	//

	// Non-blocking ch read that will timeout... boom!
	boom := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)

	for {

		// Wait for button or timeout
		log.Println("dsp.com.main: wait on a button to be pushed or a timeout")

		select {
		case <-requestBtnCh:
			log.Println("dsp.com.main: requestBtn Hit!!!!")
		case <-boom.C:
			log.Printf("dsp.com.main:  Boom! heartbeat timeout\n")
		}

		log.Println("dsp.com.main: receiveStatus()")
		receiveStatus(&mb, statusCh,&content)

		log.Println("dsp.com.main: clearDisplay()")
		dsp.ClearDisplay(&display)

		log.Println("dsp.com.main: fontExamples()")
		// dsp.FontExamples(&display)
		content.DisplayContent(&display)

		log.Println("dsp.com.main: Gosched()")
		runtime.Gosched()

	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func receiveStatus(mb *umsg.MsgBroker, statusCh chan umsg.StatusMsg, content *dsp.Content) {

	var msg umsg.StatusMsg

	select {
	case msg = <-statusCh:
		log.Printf("dsp.epaper.receiveStatus: SUCCESS, msg: [%v]\n", msg)
	default:
		log.Printf("dsp.epaper.receiveStatus: no status msg found")
	}


	//DEVTODO make this more general
	//        maybe this should be in a different function
	switch msg.Key {
	case iot.GatewayMainLoopHeartbeat:
		log.Printf("dsp.epaper.receiveStatus: call SetGatewayMainLoopHeartbeatStatus()")
		content.SetGatewayMainLoopHeartbeatStatus(msg.Value)
	default:
		log.Printf("dsp.epaper.receiveStatus: Not interested in this content: %v", msg)
	}
}
