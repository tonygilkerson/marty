package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
	"tinygo.org/x/drivers/st7789"
	"github.com/tonygilkerson/marty/pkg/marty"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

///////////////////////////////////////////////////////////////////////////////
//		main
///////////////////////////////////////////////////////////////////////////////

func main() {

	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// run light
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)
	
	//
	// Setup PIR Sensor
	//
	mboxMarty := marty.New()
	var pirCh chan string
	pirCh = make(chan string)
	go eventConsumer(pirCh, mboxMarty)
	setupPIR(pirCh)

	//
	// Setup Lora
	//
	var loraRadio *sx126x.Device
	var loraSPI = machine.SPI3
	setupLora(loraRadio,loraSPI)
	

	//
	//			Main Loop
	//
	var currentStatus, lastStatus int
	for {
		// Current status changes if any count changes
		currentStatus = mboxMarty.Ctx.ArrivedCount + mboxMarty.Ctx.DepartedCount + mboxMarty.Ctx.ErrorCount + mboxMarty.Ctx.FalseAlarmCount

		if currentStatus != lastStatus {
			lastStatus = currentStatus

			msg := fmt.Sprintf("\nArrived:  %v\nDeparted: %v\nErr:      %v\nFalse:    %v\n",
				mboxMarty.Ctx.ArrivedCount,
				mboxMarty.Ctx.DepartedCount,
				mboxMarty.Ctx.ErrorCount,
				mboxMarty.Ctx.FalseAlarmCount)
			log.Printf("\n%v\n", msg)
			loraTxRx(loraRadio, []byte(msg))

		}
		time.Sleep(time.Second * 50)
		runLight(led, 1)
	}

}

// /////////////////////////////////////////////////////////////////////////////
//
//	functions
//
// /////////////////////////////////////////////////////////////////////////////

func setupLora(loraRadio *sx126x.Device, spi  machine.SPI) {

	
	// Create the driver
	loraRadio = sx126x.New(spi)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create radio controller for target
	rc := sx126x.NewRadioControl()
	loraRadio.SetRadioController(rc)

	// Detect the device
	state := loraRadio.DetectDevice()
	if !state {
		panic("sx126x not detected.")
	}

	loraConf := lora.Config{
		Freq:           lora.MHz_916_8,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPrivate,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

}

// loraTxRx will transmit the current counts then listen for a received message
func loraTxRx(loraRadio *sx126x.Device,msg []byte) {

	//
	//	Tx
	//
	log.Printf("Send TX size=%v -> %v", len(msg), string(msg))
	err := loraRadio.Tx(msg, LORA_DEFAULT_TXTIMEOUT_MS)
	if err != nil {
		log.Printf("TX Error: %v\n", err)
	}
	
	//
	//	Rx
	//

	// DEVTODO add rx when I have a need for it

	// start := time.Now()
	// log.Println("Receiving for 5 seconds")
	// for time.Since(start) < 5*time.Second {
	// 	buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
	// 	if err != nil {
	// 		log.Printf("RX Error: %v\n", err)
	// 	} else if buf != nil {
	// 		log.Printf("Packet Received: len=%v %v\n", len(buf), string(buf))
	// 	}
	// }

}

func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 50)
		led.Low()
		time.Sleep(time.Millisecond * 50)
	}
	led.Low()
}

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}

// eventConsumer will receive event from the ISRs and send them to the state machine
func eventConsumer(ch chan string, m *marty.Marty) {
	for {
		// Wait for a change in position
		event := <-ch
		// log.Printf("eventConsumer: %v\n", event)

		// DEVTODO consider making the events of type fmt.EventID so I can remove the if statement
		if event == "FarRising" {
			m.SendEvent(marty.FarRising)
		}
		if event == "FarFalling" {
			m.SendEvent(marty.FarFalling)
		}
		if event == "NearRising" {
			m.SendEvent(marty.NearRising)
		}
		if event == "NearFalling" {
			m.SendEvent(marty.NearFalling)
		}
	}
}

// Setup PIR sensor
func setupPIR(pirCh chan string) {


	const (
		pirFar  = machine.PB10
		pirNear = machine.PA9
	)

	// Arrive Sensor
	// pirFar.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirFar.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	pirFar.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- "FarRising"
		} else {
			pirCh <- "FarFalling"
		}

	})

	// Depart Sensor
	pirNear.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirNear.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- "NearRising"
			} else {
			pirCh <- "NearFalling"
		}

	})


}