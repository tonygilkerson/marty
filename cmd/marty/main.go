package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/st7789"
	"github.com/tonygilkerson/marty/pkg/marty"
)

const (
	pirFar  = machine.PB10
	pirNear = machine.PA9
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
	runLight(led, 10)
	
	//
	// Setup PIR Sensor
	//
	mboxMarty := marty.New()
	var pirCh chan string
	pirCh = make(chan string)
	go eventConsumer(pirCh, mboxMarty)
	setupPIR(pirCh)

	//
	//			Main Loop
	//
	var currentStatus, lastStatus int
	for {

		if currentStatus != lastStatus {
			lastStatus = currentStatus

			msg := fmt.Sprintf("\nArrived:  %v\nDeparted: %v\nErr:      %v\nFalse:    %v\n",
				mboxMarty.Ctx.ArrivedCount,
				mboxMarty.Ctx.DepartedCount,
				mboxMarty.Ctx.ErrorCount,
				mboxMarty.Ctx.FalseAlarmCount)
			fmt.Printf("\n%v\n", msg)

			time.Sleep(time.Second * 3)

		}
		time.Sleep(time.Millisecond * 500)
		runLight(led, 1)
	}

}

// /////////////////////////////////////////////////////////////////////////////
//
//	functions
//
// /////////////////////////////////////////////////////////////////////////////
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

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}

// eventConsumer will receive event from the ISRs and send them to the state machine
func eventConsumer(ch chan string, m *marty.Marty) {
	for {
		// Wait for a change in position
		event := <-ch
		// log.Printf("eventConsumer: %v\n", event)

		// DEVTODO consider making the events of type fmt.EventID so I can remove the if statement
		if event == "FarRising" {
			m.SendEvent(marty.FarRising)
		}
		if event == "FarFalling" {
			m.SendEvent(marty.FarFalling)
		}
		if event == "NearRising" {
			m.SendEvent(marty.NearRising)
		}
		if event == "NearFalling" {
			m.SendEvent(marty.NearFalling)
		}
	}
}

// Setup PIR sensor
func setupPIR(pirCh chan string) {

	// Arrive Sensor
	// pirFar.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirFar.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	pirFar.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- "FarRising"
		} else {
			pirCh <- "FarFalling"
		}

	})

	// Depart Sensor
	pirNear.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirNear.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {
			pirCh <- "NearRising"
			} else {
			pirCh <- "NearFalling"
		}

	})


}