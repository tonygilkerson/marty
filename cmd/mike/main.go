package main

import (
	"fmt"
	"log"
	"machine"
	"sync"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
	"tinygo.org/x/drivers/sx126x"
)

const HEARTBEAT_DURATION_SECONDS = 600
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
	var loraTxRxMutex sync.Mutex

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
	go micMonitor(&chMicRise, loraRadio, &loraTxRxMutex)
	go mailMonitor(&chMailRise, loraRadio, &loraTxRxMutex)
	go muleMonitor(&chMuleRise, loraRadio, &loraTxRxMutex)

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	for range ticker.C {

		log.Println("------------------MainLoopHeartbeat--------------------")
		road.LoraTx(loraRadio, []byte("RoadMainLoopHeartbeat"), &loraTxRxMutex)

	}

}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func micMonitor(ch *chan string, loraRadio *sx126x.Device, mutex *sync.Mutex) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("I")

		select {

		// Heard a sound
		case <-*ch:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Sound rising")
				road.LoraTx(loraRadio, []byte("HeardSound"), mutex)
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
				// road.LoraTx(loraRadio, []byte("HeardSoundHeartbeat"), mutex)
			}
		}

	}

}

func mailMonitor(ch *chan string, loraRadio *sx126x.Device, mutex *sync.Mutex) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("M")
		select {

		// saw some light
		case <-*ch:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mailbox light rising")
				road.LoraTx(loraRadio, []byte("MailboxDoorOpened"), mutex)
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
				// road.LoraTx(loraRadio, []byte("MailboxDoorOpenedHeartbeat"), mutex)
			}
			
		}

	}

}

func muleMonitor(ch *chan string, loraRadio *sx126x.Device, mutex *sync.Mutex) {
	lastEvent := time.Now()
	lastHeartbeat := time.Now()
	active := false

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		fmt.Printf("U ")

		select {

		// light
		case <-*ch:

			lastEvent = time.Now()
			if !active {
				active = true
				log.Println("Mule light rising")
				road.LoraTx(loraRadio, []byte("MuleAlarm"), mutex)
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
				// road.LoraTx(loraRadio, []byte("MuleAlarmHeartbeat"), mutex)
			}

		}

	}

}
