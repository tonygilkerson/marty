package road

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

// setupLora will setup the lora radio device
func SetupLora(spi machine.SPI) *sx126x.Device {

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

// loraTx will transmit the current counts then listen for a received message
func LoraTx(loraRadio *sx126x.Device, msg []byte) {

	log.Printf("Send TX ------------------------------> %v", string(msg))
	err := loraRadio.Tx(msg, LORA_DEFAULT_TXTIMEOUT_MS)
	if err != nil {
		log.Printf("TX Error: %v\n", err)
	}

	start := time.Now()
	log.Println("Receiving for 1 seconds")
	for time.Since(start) < 1*time.Second {
		buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
		if err != nil {
			log.Println("RX Error: ", err)
		} else if buf != nil {
			log.Println("Packet Received: len=", len(buf), string(buf))
		}
	}
	log.Println("Receiving done.")
}