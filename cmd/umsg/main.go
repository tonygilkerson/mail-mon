package main

import (
	"log"
	"machine"
	"os"
	"runtime"
	"time"

	"github.com/tonygilkerson/mbx-iot/internal/umsg"
)

// In this setup UART0 and UART1 are connect together
// so we can perform loop back tests

// Wiring:
// 	GPIO0 -> GPIO5
// 	GPIO1 -> GPIO4

func main() {

	log.Printf("[main] - Startup pause...\n")
	time.Sleep(time.Second * 3)
	log.Printf("[main] - After startup pause\n")

	/////////////////////////////////////////////////////////////////////////////
	// Pins
	/////////////////////////////////////////////////////////////////////////////
	log.Printf("[main] - Setup\n")

	uartIn := machine.UART0
	uartInTx := machine.GPIO0
	uartInRx := machine.GPIO1

	uartOut := machine.UART1
	uartOutTx := machine.GPIO4
	uartOutRx := machine.GPIO5

	/////////////////////////////////////////////////////////////////////////////
	// Broker
	/////////////////////////////////////////////////////////////////////////////

	fooCh := make(chan umsg.FooMsg)
	barCh := make(chan umsg.BarMsg)

	// mb.SetFooCh(fooCh)

	mb := umsg.NewBroker(
		"umsg",

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

	/////////////////////////////////////////////////////////////////////////////
	// Tests
	/////////////////////////////////////////////////////////////////////////////

	fooTest(&mb, fooCh)
	barTest(&mb, barCh)

	// Done
	log.Printf("[main] - **** DONE ****")
	os.Exit(0)
}


func fooTest(mb *umsg.MsgBroker, fooCh chan umsg.FooMsg){

	var fm umsg.FooMsg
	fm.Kind = "Foo"
	fm.SenderID = umsg.LOOKBACK_SENDERID
	fm.Name = "This is a foo message from loopback"

	log.Printf("[fooTest] - PublishFoo(fm)\n")
	mb.PublishFoo(fm)

	var found bool = false
	var msg umsg.FooMsg

	// Non-blocking ch read that will timeout... boom!
	boom := time.After(3000 * time.Millisecond)
	for {
		select {
		case msg = <-fooCh:
			found = true
		case <-boom:
			log.Printf("[fooTest] - Boom! timeout waiting for message\n")
			break
		default:
			runtime.Gosched()
			time.Sleep(50 * time.Millisecond)
		}

		if found {
			break
		}
	}

	log.Printf("[fooTest] - ******************************************************************\n")
	if found {
		if msg.Name == fm.Name {
			log.Printf("[fooTest] - SUCCESS, msg: [%v]\n", msg)
		} else {
			log.Printf("[fooTest] - FAIL, wrong msg: [%v]\n", msg)
		}
	} else {
		log.Printf("[fooTest] - FAIL, did not receive message.")
	}
	log.Printf("[fooTest] - ******************************************************************\n")

}

func barTest(mb *umsg.MsgBroker, barCh chan umsg.BarMsg){

	var bm umsg.BarMsg
	bm.Kind = "Bar"
	bm.SenderID = umsg.LOOKBACK_SENDERID
	bm.Name = "This is a bar message from loopback"

	log.Printf("[barTest] - PublishBar(bm)\n")
	mb.PublishBar(bm)

	var found bool = false
	var msg umsg.BarMsg

	// Non-blocking ch read that will timeout... boom!
	boom := time.After(3000 * time.Millisecond)
	for {
		select {
		case msg = <-barCh:
			found = true
		case <-boom:
			log.Printf("[fooTest] - Boom! timeout waiting for message\n")
			break
		default:
			runtime.Gosched()
			time.Sleep(50 * time.Millisecond)
		}

		if found {
			break
		}
	}

	log.Printf("[barTest] - ******************************************************************\n")
	if found {
		if msg.Name == bm.Name {
			log.Printf("[barTest] - SUCCESS, msg: [%v]\n", msg)
		} else {
			log.Printf("[barTest] - FAIL, wrong msg: [%v]\n", msg)
		}
	} else {
		log.Printf("[barTest] - FAIL, did not receive message.")
	}
	log.Printf("[barTest] - ******************************************************************\n")

}
