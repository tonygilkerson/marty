package main

import (
	"fmt"
	"log"
	"machine"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
	"tinygo.org/x/drivers/sx127x"
)

const (
	HEARTBEAT_DURATION_SECONDS = 300
	EVENT_DURATION_SECONDS     = 3
	TICKER_MS                  = 1000
	ADC_THRESHOLD              = 2_000
)

/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	//
	// run light
	//
	time.Sleep(1 * time.Second)
	led := machine.LED //GP25
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 20)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// I would hope the channel size would never be larger than ~4 so 250 is large
	chLoraTxRx := make(chan []byte, 250)

	//
	// 	Setup Lora
	//
	loraRadio := road.SetupLora(machine.SPI0)

	//
	// Init ADC
	//
	machine.InitADC() // init the machine's ADC subsystem

	//
	// Setup Mule
	//
	muleADC := machine.ADC{Pin: machine.ADC0}

	//
	// Setup mail
	//
	mailADC := machine.ADC{Pin: machine.ADC1}

	//
	// Launch go routines
	//
	go mailMonitor(&mailADC, loraRadio, &chLoraTxRx)
	go muleMonitor(&muleADC, loraRadio, &chLoraTxRx)
	go road.LoraTx(loraRadio, &chLoraTxRx)

	// DEVTODO - remove me after test
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("STARTup-%v", i)
		chLoraTxRx <- []byte(msg)
	}

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int
	for range ticker.C {

		log.Printf("------------------MainLoopHeartbeat-------------------- %v", count)
		msg := fmt.Sprintf("RoadMainLoopHeartbeat-%v", count)
		chLoraTxRx <- []byte(msg)
		count += 1
	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print("run-")
	}

}

func mailMonitor(mailADC *machine.ADC, loraRadio *sx127x.Device, chLoraTxRx *chan []byte) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("a")

		if mailADC.Get() > ADC_THRESHOLD {
			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mailbox light rising")
				*chLoraTxRx <- []byte("MailboxDoorOpened")
			}
		} else {

			if active && time.Since(lastEvent) > EVENT_DURATION_SECONDS*time.Second {
				active = false
				log.Println("Mailbox light falling")
			}

			if time.Since(lastHeartbeat) > HEARTBEAT_DURATION_SECONDS*time.Second {
				lastHeartbeat = time.Now()
				log.Println("Mailbox Heartbeat")
				*chLoraTxRx <- []byte("MailboxDoorOpenedHeartbeat")
			}
		}
	}
}

func muleMonitor(muleADC *machine.ADC, loraRadio *sx127x.Device, chLoraTxRx *chan []byte) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("U")

		if muleADC.Get() > ADC_THRESHOLD {
			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mule light rising")
				*chLoraTxRx <- []byte("MuleAlarm")
			}
		} else {

			if active && time.Since(lastEvent) > EVENT_DURATION_SECONDS*time.Second {
				active = false
				log.Println("Mule light falling")
			}

			if time.Since(lastHeartbeat) > HEARTBEAT_DURATION_SECONDS*time.Second {
				lastHeartbeat = time.Now()
				log.Println("Mule Heartbeat")
				*chLoraTxRx <- []byte("MuleAlarmHeartbeat")
			}
		}
	}
}
