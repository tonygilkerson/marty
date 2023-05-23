package main

import (
	"log"
	"machine"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
)

/////////////////////////////////////////////////////////////////////////////
//			Main
/////////////////////////////////////////////////////////////////////////////

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	chRise := make(chan string, 1)

	//
	// 	Setup Lora
	//
	loraRadio := road.SetupLora(machine.SPI3)

	//
	// Setup Mic
	//
	micDigitalPin := machine.PA9 // if lora-e5

	micDigitalPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	micDigitalPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chRise <- "rise":
		default:
		}

	})

	lastHeard := time.Now()
	lastHeartbeat := time.Now()
	activeSound := false

	// Main loop
	for {

		select {
		//
		// Heard a sound
		//
		case <-chRise:
			lastHeard = time.Now()

			if !activeSound {
				activeSound = true
				log.Println("Sound rising")
			}

		//
		// Silence
		//
		default:

			if activeSound && time.Since(lastHeard) > 5*time.Second {
				activeSound = false
				log.Println("Sound falling")
				road.LoraTx(loraRadio, []byte("HeardSound"))
			}

			if time.Since(lastHeartbeat) > 60*time.Second {
				log.Println("Heartbeat")
				road.LoraTx(loraRadio, []byte("HeardSoundHeartbeat"))
			}

			time.Sleep(time.Millisecond * 50)
		}

	}

}
