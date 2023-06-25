package road

import (
	"log"
	"machine"
	"runtime"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx127x"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

var (
	loraRadio *sx127x.Device

	SX127X_PIN_EN   = machine.GP15
	SX127X_PIN_RST  = machine.GP20
	SX127X_PIN_CS   = machine.GP17
	SX127X_PIN_DIO0 = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	SX127X_PIN_DIO1 = machine.GP22 // (GP22--G1)I don't now what this does but it seems to need to be connected
	SX127X_SPI      = machine.SPI0
)

func dioIrqHandler(machine.Pin) {
	loraRadio.HandleInterrupt()
}

// setupLora will setup the lora radio device
func SetupLora(spi *machine.SPI) *sx127x.Device {

	spi.Configure(machine.SPIConfig{
		SCK: machine.SPI0_SCK_PIN, // GP18
		SDO: machine.SPI0_SDO_PIN, // GP19 aka MOSI
		SDI: machine.SPI0_SDI_PIN, // GP16 aka MISO
	})

	SX127X_SPI.Configure(machine.SPIConfig{Frequency: 500000, Mode: 0})
	SX127X_PIN_RST.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SX127X_PIN_EN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SX127X_PIN_EN.High() // enable the radio by default
	

	loraRadio = sx127x.New(*SX127X_SPI, SX127X_PIN_RST)
	loraRadio.SetRadioController(sx127x.NewRadioControl(SX127X_PIN_CS, SX127X_PIN_DIO0, SX127X_PIN_DIO1))

	loraRadio.Reset()
	state := loraRadio.DetectDevice()
	if !state {
		panic("main: sx127x NOT FOUND !!!")
	} else {
		log.Println("sx127x found")
	}

	// Prepare for Lora Operation
	loraConf := lora.Config{
		Freq:           lora.MHz_916_8,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPrivate,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

	return loraRadio
}

// loraTx will transmit the current counts then listen for a received message
func LoraTx(loraRadio *sx127x.Device, ch *chan string) {

	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {

		//
		// If there are no messages in the channel then get out quick
		//
		if len(*ch) == 0 {
			log.Println("LoraTx channel is empty, getting out early...")
			continue
		}

		// Enable the radio
		SX127X_PIN_EN.High()

		//
		// RX
		//
		tStart := time.Now()
		log.Println("RX Start - Receiving Lora for 5 seconds")
		for time.Since(tStart) < 5*time.Second {
		// for time.Since(tStart) < 2*time.Second {
			buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
			if err != nil {
				log.Println("RX Error: ", err)
			} else if buf != nil {
				log.Println("Packet Received: ", buf)
			}
		}
		log.Println("RX End")

		//
		// TX
		//
		log.Println("TX Start")
		var batchMsg string

		// Concatenate all messages separated by \n
		eom := false //end of messages
		for {
			select {
			case msg := <-*ch:
				if len(batchMsg) > 0 {
					batchMsg = batchMsg + "|" + msg
				} else {
					batchMsg = msg
				}
			default:
				eom = true
			}

			// break out if end of messages
			if eom {
				break
			}
		}

		//
		// Now that we have consumed all the messages from the channel Tx
		//
		if len(batchMsg) > 0 {
			log.Println("TX: ", batchMsg)
			err := loraRadio.Tx([]byte(batchMsg), LORA_DEFAULT_TXTIMEOUT_MS)
			if err != nil {
				log.Println("TX Error:", err)
			}
		} else {
			log.Println("nothing to send")	
		}
		log.Println("TX End")

		// Disable the radio to save power...
		SX127X_PIN_EN.Low()

		runtime.Gosched()
	}

}
