package main

import (
	"fmt"
	"log"
	"machine"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
	"tinygo.org/x/drivers/sx126x"
)

const HEARTBEAT_DURATION_SECONDS = 300
const EVENT_DURATION_SECONDS = 3
const TICKER_MS = 2000

/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	chMicRise := make(chan string, 1)
	chMuleRise := make(chan string, 1)
	chMailRise := make(chan string, 1)
	// I would hope you the channel size would never be larger than ~4 so 250 is large
	chLoraTxRx := make(chan []byte, 250) 

	//
	// 	Setup Lora
	//
	loraRadio := road.SetupLora(machine.SPI3)

	//
	// Setup Mic
	//
	micPin := machine.PA9 // D9 if lora-e5
	micPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	micPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chMicRise <- "Rising":
		default:
		}

	})

	//
	// Setup Mailbox
	//
	mailboxPin := machine.PA0 // D0
	mailboxPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	mailboxPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chMailRise <- "Rising":
		default:
		}

	})

	//
	// Setup Mule
	//
	mulePin := machine.PB10 // D10
	mulePin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	mulePin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chMuleRise <- "Rising":
		default:
		}

	})

	//
	// Launch go routines
	//
	go micMonitor(&chMicRise, loraRadio, &chLoraTxRx)
	go mailMonitor(&chMailRise, loraRadio, &chLoraTxRx)
	go muleMonitor(&chMuleRise, loraRadio, &chLoraTxRx)
	go road.LoraTx(loraRadio,&chLoraTxRx)

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("STARTup-%v",i)
		chLoraTxRx<-[]byte(msg)	
	}

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	for range ticker.C {

		log.Println("------------------MainLoopHeartbeat--------------------")
		chLoraTxRx<-[]byte("RoadMainLoopHeartbeat")

	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func micMonitor(chEvent *chan string, loraRadio *sx126x.Device, chLoraTxRx *chan []byte) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("I")

		select {

		// Heard a sound
		case <-*chEvent:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Sound rising")
				*chLoraTxRx<-[]byte("HeardSound")
			}

		// Silence
		default:

			if active && time.Since(lastEvent) > EVENT_DURATION_SECONDS*time.Second {
				active = false
				log.Println("Sound falling")
			}

			if time.Since(lastHeartbeat) > HEARTBEAT_DURATION_SECONDS*time.Second {
				lastHeartbeat = time.Now()
				log.Println("Mic Heartbeat")
				// *chLoraTxRx<-[]byte("HeardSoundHeartbeat")
			}
		}

	}

}

func mailMonitor(chEvent *chan string, loraRadio *sx126x.Device, chLoraTxRx *chan []byte) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("M")
		select {

		// saw some light
		case <-*chEvent:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mailbox light rising")
				*chLoraTxRx<-[]byte("MailboxDoorOpened")
			}

		// dark
		default:

			if active && time.Since(lastEvent) > EVENT_DURATION_SECONDS*time.Second {
				active = false
				log.Println("Mailbox light falling")
			}

			if time.Since(lastHeartbeat) > HEARTBEAT_DURATION_SECONDS*time.Second {
				lastHeartbeat = time.Now()
				log.Println("Mailbox Heartbeat")
				//*chLoraTxRx<-[]byte("MailboxDoorOpenedHeartbeat")
			}
			
		}

	}

}

func muleMonitor(chEvent *chan string, loraRadio *sx126x.Device, chLoraTxRx *chan []byte) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("U ")

		select {

		// light
		case <-*chEvent:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mule light rising")
				*chLoraTxRx<-[]byte("MuleAlarm")
			}

		// dark
		default:

			if active && time.Since(lastEvent) > EVENT_DURATION_SECONDS*time.Second {
				active = false
				log.Println("Mule light falling")
			}

			if time.Since(lastHeartbeat) > HEARTBEAT_DURATION_SECONDS*time.Second {
				lastHeartbeat = time.Now()
				log.Println("Mule Heartbeat")
				// *chLoraTxRx<-[]byte("MuleAlarmHeartbeat")
			}

		}

	}

}
