package main

import (
	"log"
	"machine"
	"math"
	"runtime"
	"time"

	tm1637mod "github.com/tonygilkerson/mbx-iot/hack/driver/tm1637"
	"github.com/tonygilkerson/mbx-iot/internal/dsp"
	"github.com/tonygilkerson/mbx-iot/internal/road"
	"github.com/tonygilkerson/mbx-iot/internal/soil"
	"github.com/tonygilkerson/mbx-iot/internal/util"
	"github.com/tonygilkerson/mbx-iot/pkg/iot"

	"tinygo.org/x/drivers/sx127x"

)

func main() {

	//
	// Named PINs
	//
	var tm1637CLK machine.Pin = machine.GP10
	var tm1637DIO machine.Pin = machine.GP11
	var soilSDA machine.Pin = machine.GP12
	var soilSCL machine.Pin = machine.GP13
	var loraEn machine.Pin = machine.GP15
	var loraSdi machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var loraCs machine.Pin = machine.GP17
	var loraSck machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var loraSdo machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN
	var loraRst machine.Pin = machine.GP20
	var loraDio0 machine.Pin = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	var loraDio1 machine.Pin = machine.GP22 // (GP22--G1) I don't now what this does but it seems to need to be connected

	var led machine.Pin = machine.GPIO25 // GP25 machine.LED

	const (
		HEARTBEAT_DURATION_SECONDS        = 15
		TXRX_LOOP_TICKER_DURATION_SECONDS = 10
	)

	//
	// run light
	//
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsp.RunLight(led, 10)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// Configure 4 digit 7-segment display
	//
	tm := tm1637mod.New(tm1637CLK, tm1637DIO, 5)

	//
	// Configure I2C
	//
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{SDA: soilSDA, SCL: soilSCL})
	util.DoOrDie(err)

	soil := soil.New(i2c)

	//
	// 	Setup Lora
	//
	var loraRadio *sx127x.Device
	txQ := make(chan string, 250) // I would hope the channel size would never be larger than ~4 so 250 is large
	rxQ := make(chan string, 250)

	log.Println("Setup LORA")
	radio := road.SetupLora(*machine.SPI0, loraEn, loraRst, loraCs, loraDio0, loraDio1, loraSck, loraSdo, loraSdi, loraRadio, &txQ, &rxQ, 0, 10_000, TXRX_LOOP_TICKER_DURATION_SECONDS, road.TxRx)

	// Routine to send and receive
	go radio.LoraRxTxRunner()

	//
	// Main loop
	//
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int

	for range ticker.C {

		log.Printf("------------------SoilMainLoopHeartbeat-------------------- %v", count)
		count += 1

		//
		// Send Heartbeat to Tx queue
		//
		txQ <- iot.SoilMainLoopHeartbeat
		dsp.RunLight(led, 2)

		m, err := soil.ReadMoisture()
		util.DoOrDie(err)
		log.Printf("Moisture: %v\n", m)
		time.Sleep(time.Second)

		t, err := soil.ReadTemperature()
		util.DoOrDie(err)
		log.Printf("Temperature (F): %v\n", t)
		time.Sleep(time.Second)

		// alternate between displaying the moisture and temperature
		if math.Mod(float64(count), 2) == 0 {
			tm.DisplayNumber(int16(m))
		} else {
			tm.DisplayNumber(int16(t))
		}

		//
		// Let someone else have a turn
		//
		runtime.Gosched()
	}
}
