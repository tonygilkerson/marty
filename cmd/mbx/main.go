package main

import (
	"device/arm"
	"log"
	"machine"
	"time"

	"github.com/tonygilkerson/marty/pkg/fsm"
	"github.com/tonygilkerson/marty/pkg/marty"

	// import "device/arm"
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
	log.Printf("Starting...")
	

	//
	// 	Setup PIR Sensor and start the event consumer
	//
	mbx := marty.New()
	var pirCh chan fsm.EventID
	pirCh = make(chan fsm.EventID)
	go eventConsumer(pirCh, mbx)
	setupPIR(pirCh)

	//
	// 	Setup Lora
	//
	loraRadio := setupLora(machine.SPI3)

	//
	//	Publish Metrics
	//
	go publishMetrics(mbx, loraRadio, led)

	// Reset device every so often
	for {

		time.Sleep(time.Hour * 12)
		// runLight(led, 30)
		// log.Printf("SystemReset...")
		arm.SystemReset()
	}
}

///////////////////////////////////////////////////////////////////////////////
//															functions
//////////////////////////////////////////////////////////////////////////////

// publishMetrics will publish the mbox status via Lora on a schedule
func publishMetrics(mbx *marty.Marty, loraRadio *sx126x.Device, led machine.Pin) {

	var lastArrivedCount, lastDepartedCount, lastErrorCount, lastFalseAlarmCount, loopCount int
	
	for {
		
		if mbx.Ctx.ArrivedCount != lastArrivedCount {
			lastArrivedCount = mbx.Ctx.ArrivedCount
			log.Printf("Tx: Arrived")
			loraTx(loraRadio, []byte(marty.Arrived))
			runLight(led, 2)
		}

		if mbx.Ctx.DepartedCount != lastDepartedCount {
			lastDepartedCount = mbx.Ctx.DepartedCount
			log.Printf("Tx: Departed")
			loraTx(loraRadio, []byte(marty.Departed))
			runLight(led, 2)
		}

		if mbx.Ctx.ErrorCount != lastErrorCount {
			lastErrorCount = mbx.Ctx.ErrorCount
			log.Printf("Tx: Error")
			loraTx(loraRadio, []byte(marty.Error))
			runLight(led, 2)
		}

		if mbx.Ctx.FalseAlarmCount != lastFalseAlarmCount {
			lastFalseAlarmCount = mbx.Ctx.FalseAlarmCount
			log.Printf("Tx: FalseAlarm")
			loraTx(loraRadio, []byte(marty.FalseAlarm))
			runLight(led, 2)
		}
		
		// I am not sure what the best delay should be here but if it is too large
		// multiple Arrivals for example will only get counted as one
		time.Sleep(time.Second * 5)
		
		// Send a heartbeat every minute
		loopCount += 1
		if loopCount > 12 {
			loopCount = 0
			loraTx(loraRadio, []byte("MBX-HEARTBEAT"))
			runLight(led, 2)
		}
		
		
	}

}

// loraTx will transmit the current counts then listen for a received message
func loraTx(loraRadio *sx126x.Device, msg []byte) {

	log.Printf("Send TX size=%v -> %v", len(msg), string(msg))
	err := loraRadio.Tx(msg, LORA_DEFAULT_TXTIMEOUT_MS)
	if err != nil {
		log.Printf("TX Error: %v\n", err)
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
		led.Low()
		time.Sleep(time.Millisecond * 200)
		// Do high last because we want it to be off and for some reason
		// high is off on lore E5 board, strange
		led.High()
		time.Sleep(time.Millisecond * 200)
	}

}

// eventConsumer will receive event from the ISRs and send them to the state machine
func eventConsumer(ch chan fsm.EventID, m *marty.Marty) {

	var event fsm.EventID
	for {
		// Wait for a change in position
		event = <-ch

		if err := m.StateMachine.SendEvent(event, &m.Ctx); err == fsm.ErrEventRejected {
			// log.Printf("Error: %v\n", event)
			m.Ctx.ErrorCount += 1
			m.StateMachine.Current = fsm.Default
		}
		
	}
}

// Setup PIR sensor
func setupPIR(pirCh chan fsm.EventID) {

	const (
		pirNear = machine.PB10
		pirFar  = machine.PA9
	)

	// Arrive Sensor
	// pirFar.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirFar.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	pirFar.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- marty.FarRising
		} else {
			pirCh <- marty.FarFalling
		}

	})

	// Depart Sensor
	pirNear.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirNear.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- marty.NearRising
		} else {
			pirCh <- marty.NearFalling
		}

	})

}
