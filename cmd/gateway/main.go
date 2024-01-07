package main

import (
	"fmt"
	"log"
	"machine"
	"runtime"
	"strings"
	"time"

	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"github.com/tonygilkerson/mbx-iot/internal/road"
	"github.com/tonygilkerson/mbx-iot/pkg/iot"
	"tinygo.org/x/drivers/sx127x"
)

const (
	HEARTBEAT_DURATION_SECONDS        = 8
	TXRX_LOOP_TICKER_DURATION_SECONDS = 6
)

/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	//
	// Named PINs
	//
	var en machine.Pin = machine.GP15
	var sdi machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var cs machine.Pin = machine.GP17
	var sck machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var sdo machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN
	var rst machine.Pin = machine.GP20
	var dio0 machine.Pin = machine.GP21  // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	var dio1 machine.Pin = machine.GP22  // (GP22--G1)I don't now what this does but it seems to need to be connected
	var uartTx machine.Pin = machine.GP0 // machine.UART0_TX_PIN
	var uartRx machine.Pin = machine.GP1 // machine.UART0_RX_PIN
	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsp.RunLight(led, 10)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// setup Uart
	//
	log.Println("Configure UART")
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{BaudRate: 115200, TX: uartTx, RX: uartRx})

	//
	// 	Setup Lora
	//
	var loraRadio *sx127x.Device
	// I am thinking that a batch of message can be half dozen max so 250 should be plenty large
	txQ := make(chan string, 250)
	rxQ := make(chan string, 250)

	log.Println("Setup LORA")
	radio := road.SetupLora(*machine.SPI0, en, rst, cs, dio0, dio1, sck, sdo, sdi, loraRadio, &txQ, &rxQ, 0, 10_000, TXRX_LOOP_TICKER_DURATION_SECONDS, road.TxRx)

	// Create status map
	status := make(map[string]string)

	// Launch go routines
	log.Println("Launch go routines")
	go writeToSerial(&rxQ, uart, status)
	go readFromSerial(&txQ, uart)
	go radio.LoraRxTxRunner()

	// Main loop
	log.Println("Start main loop")

	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int

	for range ticker.C {

		log.Printf("------------------mbx-iot gateway MainLoopHeartbeat-------------------- %v", count)
		count += 1
		updateStatus(iot.GatewayMainLoopHeartbeat, status)

		//
		// Send out status on each heartbeat
		//
		for k, v := range status {
			txQ <- fmt.Sprintf("%s:%s", k, v)
		}

		dsp.RunLight(led, 2)
		runtime.Gosched()
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func writeToSerial(rxQ *chan string, uart *machine.UART, status map[string]string) {
	var msgBatch string

	for msgBatch = range *rxQ {

		log.Printf("gateway.writeToSerial: Message batch: [%v]", msgBatch)

		messages := road.SplitMessageBatch(msgBatch)
		for _, msg := range messages {
			log.Printf("gateway.writeToSerial: Write to serial: [%v]", msg)
			uart.Write([]byte(msg))
			updateStatus(msg, status)
			time.Sleep(time.Millisecond * 50) // Mark the End of a message
		}

		runtime.Gosched()

	}

}

//
// readFromSerial will read messages sent from the cluster and broadcast them for receive in the field
//                currently this is used for testing. The cluster exposed a REST endpoint that can be
//                post a message that is subsequently read and then transmitted here
//
func readFromSerial(txQ *chan string, uart *machine.UART) {
	data := make([]byte, 250)

	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {

		//
		// Check to see if we have any data to read
		//
		if uart.Buffered() == 0 {
			//Serial buffer is empty, nothing to do, get out!"
			continue
		}

		//
		// Read from serial then transmit the message
		//
		n, err := uart.Read(data)
		if err != nil {
			log.Printf("Serial read error [%v]", err)
		} else {
			log.Printf("Put on txQ [%v]", string(data[:n]))
			*txQ <- string(data[:n])
		}

		runtime.Gosched()
	}

}

func updateStatus(msg string, status map[string]string) {
	// Get now as a string
	t := time.Now()
	ts := fmt.Sprintf(t.Format("20060102150405"))

	//
	// Each message is a a key:values pair
	//
	parts := strings.Split(msg, ":")
	var msgKey string
	var msgValue string

	if len(parts) > 0 {
		msgKey = parts[0]
	}
	if len(parts) > 1 {
		msgValue = parts[1]
	}


	// Update status for given
	switch {

	case msgKey == string(iot.MbxTemperature):
		status[msgKey] = msgValue
		log.Printf("gateway.updateStatus: update %s with status: %s\n", msgKey, msgValue)
		
	default:
		log.Printf("gateway.updateStatus: update %s with status: %s\n", msgKey, ts)
		status[msgKey] = ts
	}

}
