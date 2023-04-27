package main

import (
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
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
	runLight(led, 5)

	//
	// setupUart
	//
	uart := machine.UART2
	machine.UART2.Configure(machine.UARTConfig{BaudRate: 9600, TX: machine.UART2_TX_PIN, RX: machine.UART2_RX_PIN})

	//
	// 	Setup Lora
	//
	// loraRadio := setupLora(machine.SPI3)

	//
	// Monitor LoraRx
	//
	// go loraRx(uart, loraRadio)

	//
	// Send heartbeat every minute 
	//
	for {

		runLight(led, 2)
		log.Printf("HEARTBEAT")
		uart.Write([]byte("HEARTBEAT"))
		time.Sleep(time.Second * 60)

	}

}

///////////////////////////////////////////////////////////////////////////////
//															functions
//////////////////////////////////////////////////////////////////////////////


// loraRx will receive messages from in-the-field IOT devices
func loraRx(uart *machine.UART,  loraRadio *sx126x.Device)  {

	for {
		buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)

		if err != nil {
			log.Printf("RX Error: %v", err)
		}

		if buf != nil {
			log.Printf("Packet received, send to UART: len: %v, msg: %v", len(buf), string(buf))
			uart.Write(buf)
		} else {
			log.Printf("No packets to receive")
		}
	}

}

// setupLora will setup the lora radio device
func setupLora(spi machine.SPI) *sx126x.Device {

	var loraRadio *sx126x.Device

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

	return loraRadio
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
