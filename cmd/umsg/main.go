package main

import (
	"fmt"
	"log"
	"machine"
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
	// Foo test
	/////////////////////////////////////////////////////////////////////////////
	
	log.Printf("[main] - Buildup foo message 1\n")
	var fm umsg.FooMsg
	fm.Kind = "Foo"
	fm.SenderID = "" // will default to broker senderID
	fm.Name = "This is a foo message"

	
	log.Printf("[main] - PublishFoo(fm)\n")
	mb.PublishFoo(fm)
	
	log.Printf("[main] - Gosched()\n")
	runtime.Gosched()
	
	log.Printf("[main] - pause after publish to give the listenRoutine() time\n")
	time.Sleep(time.Millisecond * 4000)

	wantChLen := 0
	gotChLen := len(fooCh)

	log.Printf("[main] - ******************************************************************\n")
	if wantChLen == gotChLen {
		log.Printf("[main] - SUCCESS\n")
	} else {
		log.Printf("[main] - FAIL, want: [%v], got: [%v]\n", wantChLen, gotChLen)
	}
	log.Printf("[main] - ******************************************************************\n")

	/////////////////////////////////////////////////////////////////////////////
	// Foo test 2
	/////////////////////////////////////////////////////////////////////////////
	
	log.Printf("[main] - pause between tests **************************************************************** \n")
	time.Sleep(time.Millisecond * 5000)

	log.Printf("[main] - Buildup foo message from loopback\n")
	// var fm umsg.FooMsg
	fm.Kind = "Foo"
	fm.SenderID = umsg.LOOKBACK_SENDERID
	fm.Name = "This is a foo message from loopback"

	
	
	log.Printf("[main] - PublishFoo(fm)\n")
	mb.PublishFoo(fm)
	
	log.Printf("[main] - Gosched()\n")
	runtime.Gosched()

	log.Printf("[main] - pause after publish to give the listenRoutine() time\n")
	time.Sleep(time.Millisecond * 4000)

	// for i := range fooCh {
	// 	log.Printf("[main] - DEBUG GOT , msg: ********************************** [%v] *********************************************\n",i)
	// }

	var found bool = false 
	var msg umsg.FooMsg

	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(1000 * time.Millisecond)
	for {
		select {
		case msg = <-fooCh:
			found = true
			log.Printf("[main] - got message: %v\n", msg)
		case <-tick:
		case <-boom:
			break
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}

		if found {
			break
		}
	}


	log.Printf("[main] - ******************************************************************\n")
	if  found {
		log.Printf("[main] - SUCCESS, msg: [%v]\n",msg)
	} else {
		log.Printf("[main] - FAIL, want: [%v], got: [%v]\n", wantChLen, gotChLen)
	}
	log.Printf("[main] - ******************************************************************\n")

	/////////////////////////////////////////////////////////////////////////////
	// Keep Alive...
	/////////////////////////////////////////////////////////////////////////////

	for {
		fmt.Printf("m")
		runtime.Gosched()
		time.Sleep(time.Millisecond * 500)
	}
}
