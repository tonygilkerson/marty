package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"

	// "tinygo.org/x/tinyfont"
	// "tinygo.org/x/tinyfont/freemono"

	"github.com/tonygilkerson/marty/pkg/marty"
)

const (
	pirArrive = machine.GP20
	pirDepart = machine.GP21
)

///////////////////////////////////////////////////////////////////////////////
//		main
///////////////////////////////////////////////////////////////////////////////

func main() {
	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// run light
	//
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)

	// Create a new instance of the state machine.
	mboxMarty := marty.New()

	var pirCh chan string
	pirCh = make(chan string)
	go eventConsumer(pirCh, &mboxMarty)

	// Arrive Sensor
	pirArrive.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirArrive.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {

			fmt.Printf("ISR PIR Arriving ArriveRising PinRising\n")
			pirCh <- "ArriveRising"

		} else {

			fmt.Printf("ISR PIR Arriving ArriveFalling PinFalling\n")
			pirCh <- "ArriveFalling"

		}

	})

	// Depart Sensor
	pirDepart.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirDepart.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

		if p.Get() {

			fmt.Printf("ISR PIR Departing DepartRising PinRising\n")
			pirCh <- "DepartRising"

		} else {

			fmt.Printf("ISR PIR Departing DepartFalling PinFalling\n")
			pirCh <- "DepartFalling"

		}

	})

	//
	// setup the display
	//
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  0,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP28, // I don't think this is actually used for LCD, just assign to any open pin
	})

	display := st7789.New(machine.SPI1,
		machine.GP12, // TFT_RESET
		machine.GP8,  // TFT_DC
		machine.GP9,  // TFT_CS
		machine.GP13) // TFT_LITE

	display.Configure(st7789.Config{
		// With the display in portrait and the usb socket on the left and in the back
		// the actual width and height are switched width=320 and height=240
		Width:        240,
		Height:       320,
		Rotation:     st7789.ROTATION_90,
		RowOffset:    0,
		ColumnOffset: 0,
		FrameRate:    st7789.FRAMERATE_111,
		VSyncLines:   st7789.MAX_VSYNC_SCANLINES,
	})

	width, height := display.Size()
	fmt.Printf("width: %v, height: %v\n", width, height)

	// red := color.RGBA{126, 0, 0, 255} // dim
	// red := color.RGBA{255, 0, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	green := color.RGBA{0, 255, 0, 255}

	//
	// Setup input buttons (the ones on the display)
	//

	// If any key is pressed record the corresponding pin
	var keyPressed machine.Pin

	key0 := machine.GP15
	key1 := machine.GP17
	key2 := machine.GP2
	key3 := machine.GP3

	key0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	key0.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
		fmt.Printf("Key0\n")
	})
	key1.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
		fmt.Printf("Key1\n")
	})
	key2.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
		fmt.Printf("Key2\n")
	})
	key3.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
		fmt.Printf("Key3\n")
	})

	//
	//			Main Loop
	//
	var currentStatus, lastStatus int
	for {

		if keyPressed == key0 {
			keyPressed = 0
			log.Printf("key0\n")
		}
		if keyPressed == key1 {
			keyPressed = 0
			log.Printf("key1\n")
			mboxMarty.ResetContext()
		}
		if keyPressed == key2 {
			keyPressed = 0
			log.Printf("key2\n")
		}
		if keyPressed == key3 {
			keyPressed = 0
			log.Printf("key3\n")
		}

		//
		// Display
		//
		currentStatus = mboxMarty.Ctx.ArrivedCount + mboxMarty.Ctx.DepartedCount + mboxMarty.Ctx.ErrorCount + mboxMarty.Ctx.FalseAlarmCount

		if currentStatus != lastStatus {
			lastStatus = currentStatus

			cls(&display)
			msg := fmt.Sprintf("Arrived:  %v\nDeparted: %v\nErr:      %v\nFalse:    %v",
				mboxMarty.Ctx.ArrivedCount,
				mboxMarty.Ctx.DepartedCount,
				mboxMarty.Ctx.ErrorCount,
				mboxMarty.Ctx.FalseAlarmCount)

			tinyfont.WriteLine(&display, &freemono.Regular18pt7b, 10, 30, msg, green)
			fmt.Printf("%v\n\n", msg)
		}
		time.Sleep(time.Millisecond * 200)
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
		// fmt.Printf("event: %v\n", event)

		// DEVTODO consider making the events of type fmt.EventID so I can remove the if statement
		if event == "ArriveRising" {
			fmt.Printf("SEND: %v\n", marty.ArriveRising)
			m.SendEvent(marty.ArriveRising)
		}
		if event == "ArriveFalling" {
			fmt.Printf("SEND: %v\n", marty.ArriveFalling)
			m.SendEvent(marty.ArriveFalling)
		}
		if event == "DepartRising" {
			fmt.Printf("SEND: %v\n", marty.DepartRising)
			m.SendEvent(marty.DepartRising)
		}
		if event == "DepartFalling" {
			m.SendEvent(marty.DepartFalling)
			fmt.Printf("SEND: %v\n", marty.DepartFalling)
		}
	}
}
