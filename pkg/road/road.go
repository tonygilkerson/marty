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
func LoraTx(loraRadio *sx126x.Device, ch *chan []byte) {

	for msg := range *ch {

		log.Printf("Send TX ------------------------------> %v", string(msg))
		err := loraRadio.Tx(msg, LORA_DEFAULT_TXTIMEOUT_MS)
		if err != nil {
			log.Printf("TX Error: %v, sending msg: %v\n", err, string(msg))
		}

		start := time.Now()
		log.Printf("Receiving for up to 10 seconds after msg: %v", string(msg))
		for time.Since(start) < 10*time.Second {
			log.Printf("loraRadio.Rx...\n")
			buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
			if err != nil {
				log.Printf("RX Error: %v, after msg: %v", err, string(msg))
			}
			if buf != nil {
				log.Printf("<----------Packet Received: %v, after msg: %v", string(buf), string(msg))
				break
			}

		}
		log.Printf("Receiving done after msg: %v", string(msg))
		//DEVTODO not sure if this is needed but I feel like we need to wait just a bit before trying to send again
		//        to give the receiver time to do its thing and start listening agin
		time.Sleep(time.Millisecond * 5000)

	}

}
