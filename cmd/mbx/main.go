package main

import (
	"fmt"
	"log"
	"machine"
	"runtime"
	"time"

	"github.com/tonygilkerson/marty/pkg/road"
)

const (
	HEARTBEAT_DURATION_SECONDS = 300
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
	chLoraTxRx := make(chan string, 250)

	//
	// 	Setup Lora
	//
	loraRadio := road.SetupLora(machine.SPI0)

	//
	// Setup charger
	//

	// CHG - Charge status (active low) pulls to GND (open drain) lighting the connected led when the battery is charging.
	// If the battery is charged or the charger is disabled, CHG is disconnected from ground (high impedance) and the LED will be off.
	chgPin := machine.GP10
	chgPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	log.Printf("chgPin status: %v\n", chgPin.Get())

	// PGOOD - Power Good Status (active low). PGOOD pulls to GND (open drain) lighting the connected led when a valid input source is connected.
	// If the input power source is not within specified limits, PGOOD is disconnected from ground (high impedance) and the LED will be off.
	pgoodPin := machine.GP11
	pgoodPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	log.Printf("pgoodPin status: %v\n", pgoodPin.Get())

	//
	// Setup Mule
	//
	muleCh := make(chan string)
	mulePin := machine.GP12
	mulePin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	log.Printf("mulePin status: %v\n", mulePin.Get())

	mulePin.SetInterrupt(machine.PinRising, func(p machine.Pin) {

		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case muleCh <- "up":
		default:
		}

	})

	//
	// Setup mail
	//
	mailCh := make(chan string)
	mailPin := machine.GP13
	mailPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	log.Printf("mailPin status: %v\n", mulePin.Get())

	mailPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {

		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case mailCh <- "up":
		default:
		}

	})

	// Launch go routines
	go mailMonitor(&mailCh, &chLoraTxRx)
	go muleMonitor(&muleCh, &chLoraTxRx)
	go road.LoraTx(loraRadio, &chLoraTxRx)

	// Main loop
	ticker := time.NewTicker(time.Second * HEARTBEAT_DURATION_SECONDS)
	var count int
	for range ticker.C {

		log.Printf("------------------MainLoopHeartbeat-------------------- %v", count)
		count += 1
		log.Printf("mailPin status: %v\n", mailPin.Get())
		log.Printf("mulePin status: %v\n", mulePin.Get())

		//
		// Send Heartbeat to Tx queue
		//
		chLoraTxRx <- "RoadMainLoopHeartbeat"

		//
		// send charger status
		//
		sendChargerStatus(chgPin, pgoodPin, &chLoraTxRx)

		//
		// Send Temperature to Tx queue
		//
		sendTemperature(&chLoraTxRx)

		runtime.Gosched()
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

func mailMonitor(ch *chan string, chLoraTxRx *chan string) {

	for range *ch {
		log.Println("Mailbox light up")
		*chLoraTxRx <- "MailboxDoorOpened"

		runtime.Gosched()
		// Wait a long time to give mail man time to shut the door
		time.Sleep(time.Second * 60)
		log.Println("Mailbox light down")

	}

}

func muleMonitor(ch *chan string, chLoraTxRx *chan string) {

	for range *ch {
		log.Println("Mule light up")
		*chLoraTxRx <- "MuleAlarm"

		runtime.Gosched()
		time.Sleep(time.Second * 4)
		log.Println("Mule light down")

	}
}

func sendTemperature(chLoraTxRx *chan string) {

	// F = ( (ReadTemperature /1000) * 9/5) + 32
	fahrenheit := ((machine.ReadTemperature() / 1000) * 9 / 5) + 32
	fmt.Printf("fahrenheit: %v\n", fahrenheit)
	*chLoraTxRx <- fmt.Sprintf("MailboxTemperature:%v", fahrenheit)

}

func sendChargerStatus(chgPin machine.Pin, pgoodPin machine.Pin, chLoraTxRx *chan string) {

	if pgoodPin.Get() {
		log.Println("Power source bad")
		*chLoraTxRx <- "ChargerPowerSourceBad"
	} else {
		log.Println("Power source good")
		*chLoraTxRx <- "ChargerPowerSourceGood"
	}

	if chgPin.Get() {
		log.Println("Charger off")
		*chLoraTxRx <- "ChargerChargeStatusOff"
	} else {
		log.Println("Charger on")
		*chLoraTxRx <- "ChargerChargeStatusOn"
	}

}
