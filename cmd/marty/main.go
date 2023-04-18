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
	pirFar  = machine.GP20
	pirNear = machine.GP21
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

	//
	// Create channel for PIR events
	//
	var pirCh chan string
	pirCh = make(chan string)
	go eventConsumer(pirCh, mboxMarty)

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

	//
	// Setup Display
	//
	display, key0, key1, key2, key3, displayKeyPressed := setupDisplay()

	//
	//			Main Loop
	//
	var currentStatus, lastStatus int
	for {

		if *displayKeyPressed == key0 {
			*displayKeyPressed = 0
			log.Printf("key0\n")
		}
		if *displayKeyPressed == key1 {
			*displayKeyPressed = 0
			log.Printf("key1\n")
			mboxMarty.ResetContext()
		}
		if *displayKeyPressed == key2 {
			*displayKeyPressed = 0
			log.Printf("key2\n")
		}
		if *displayKeyPressed == key3 {
			*displayKeyPressed = 0
			log.Printf("key3\n")
		}

		// red := color.RGBA{126, 0, 0, 255} // dim
		// red := color.RGBA{255, 0, 0, 255}
		// black := color.RGBA{0, 0, 0, 255}
		// white := color.RGBA{255, 255, 255, 255}
		// blue := color.RGBA{0, 0, 255, 255}
		green := color.RGBA{0, 255, 0, 255}

		//
		// Display
		//
		currentStatus = mboxMarty.Ctx.ArrivedCount + mboxMarty.Ctx.DepartedCount + mboxMarty.Ctx.ErrorCount + mboxMarty.Ctx.FalseAlarmCount

		if currentStatus != lastStatus {
			lastStatus = currentStatus

			cls(&display)
			display.EnableBacklight(true)
			msg := fmt.Sprintf("\nArrived:  %v\nDeparted: %v\nErr:      %v\nFalse:    %v\n",
				mboxMarty.Ctx.ArrivedCount,
				mboxMarty.Ctx.DepartedCount,
				mboxMarty.Ctx.ErrorCount,
				mboxMarty.Ctx.FalseAlarmCount)

			tinyfont.WriteLine(&display, &freemono.Regular18pt7b, 10, 30, msg, green)
			fmt.Printf("\n%v\n", msg)

			time.Sleep(time.Second * 3)
			display.EnableBacklight(false)
			cls(&display)
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

func setupDisplay() (display st7789.Device, key0, key1, key2, key3  machine.Pin, displayKeyPressed *machine.Pin) {

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

	display = st7789.New(machine.SPI1,
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

	//
	// Setup input buttons (the ones on the display)
	//
	key0 = machine.GP15
	key1 = machine.GP17
	key2 = machine.GP2
	key3 = machine.GP3

	key0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	key0.SetInterrupt(machine.PinFalling, func(p machine.Pin) { displayKeyPressed = &p })
	key1.SetInterrupt(machine.PinFalling, func(p machine.Pin) { displayKeyPressed = &p })
	key2.SetInterrupt(machine.PinFalling, func(p machine.Pin) { displayKeyPressed = &p })
	key3.SetInterrupt(machine.PinFalling, func(p machine.Pin) { displayKeyPressed = &p })

	width, height := display.Size()
	fmt.Printf("width: %v, height: %v\n", width, height)

	// Off and clear initially
	display.EnableBacklight(false)
	cls(&display)


	return
	
}