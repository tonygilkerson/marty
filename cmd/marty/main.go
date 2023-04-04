package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/st7789"
	// "tinygo.org/x/tinyfont"
	// "tinygo.org/x/tinyfont/freemono"

	"github.com/tonygilkerson/marty/pkg/fsm"
)

const (
	// States
	// Default    fsm.StateID = "Default"
	Arriving   fsm.StateID = "Arriving"
	Arrived    fsm.StateID = "Arrived"
	Departing  fsm.StateID = "Departing"
	Departed   fsm.StateID = "Departed"
	FalseAlarm fsm.StateID = "FalseAlarm"
	Error      fsm.StateID = "Error"

	//Events
	RightRising  fsm.EventID = "RightRising"
	RightFalling fsm.EventID = "RightFalling"
	LeftRising   fsm.EventID = "LeftRising"
	LeftFalling  fsm.EventID = "LeftFalling"
	Reset        fsm.EventID = "Reset"
)

type MartyContext struct {
	DefaultCount    int
	ArrivedCount    int
	ArrivingCount   int
	DepartedCount   int
	DepartingCount  int
	ErrorCount      int
	FalseAlarmCount int
}

func (c *MartyContext) String() string {
	cCopy := *c
	return fmt.Sprintf("MartyContext: %+v\n", cCopy)
}

// DefaultAction
type DefaultAction struct{}

func (a *DefaultAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DefaultCount += 1

	log.Printf("DefaultAction\n\n")
	return fsm.NoOp
}

// ArrivedAction
type ArrivedAction struct{}

func (a *ArrivedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ArrivedCount += 1

	log.Printf("ArrivedAction\n")
	return Reset
}

// ArrivingAction
type ArrivingAction struct{}

func (a *ArrivingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ArrivingCount += 1

	log.Printf("ArrivingAction\n")
	return fsm.NoOp
}

// DepartedAction
type DepartedAction struct{}

func (a *DepartedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DepartedCount += 1

	log.Printf("DepartedAction\n")
	return Reset
}

// DepartingAction
type DepartingAction struct{}

func (a *DepartingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DepartingCount += 1

	log.Printf("DepartingAction\n")
	return fsm.NoOp
}

// ErrorAction
type ErrorAction struct{}

func (a *ErrorAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ErrorCount += 1

	log.Printf("ErrorAction\n")
	return fsm.NoOp
}

// FalseAlarmAction
type FalseAlarmAction struct{}

func (a *FalseAlarmAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.FalseAlarmCount += 1

	log.Printf("FalseAlarmAction\n")
	return Reset
}

func newMartyFSM() *fsm.StateMachine {

	fsm := fsm.StateMachine{
		Current:  fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{

			fsm.Default: fsm.State{
				Action: &DefaultAction{},
				Events: fsm.Events{
					RightRising: Arriving,
					LeftRising:  Departing,
				},
			},

			Arriving: fsm.State{
				Action: &ArrivingAction{},
				Events: fsm.Events{
					LeftRising:   Arrived,
					RightFalling: FalseAlarm,
				},
			},

			Arrived: fsm.State{
				Action: &ArrivedAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},

			Departing: fsm.State{
				Action: &DepartingAction{},
				Events: fsm.Events{
					LeftFalling: FalseAlarm,
					RightRising: Departed,
				},
			},

			Departed: fsm.State{
				Action: &DepartedAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},

			FalseAlarm: fsm.State{
				Action: &FalseAlarmAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},
		},
	}

	// log.Printf("\n-------\n\n%+v\n\n--------\n",fsm)

	return &fsm
}

const (
	pirRight = machine.GP21
)
var pirCh chan string
var ctx MartyContext
var sm *fsm.StateMachine

func pirIsr(p machine.Pin) {

	// disable the interrupt service routine since that's a good thing to do.
	pirRight.SetInterrupt(machine.PinFalling|machine.PinRising, nil)

	if p.Get() {

		fmt.Printf("PIR arriving PinRising event\n")
		// sendEvent(RightRising, sm, &ctx)
		pirCh <- "RightRising"

	} else {

		fmt.Printf("PIR arriving PinFalling event\n")
		// sendEvent(RightFalling, sm, &ctx)
		pirCh <- "RightFalling"
	}

	// re-enable the interrupt service routine.
	pirRight.SetInterrupt(machine.PinFalling|machine.PinRising, pirIsr)

}

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
	sm = newMartyFSM()
	pirCh = make(chan string)
	pirRight.Configure(machine.PinConfig{Mode: machine.PinInput})
	pirRight.SetInterrupt(machine.PinFalling|machine.PinRising, pirIsr)
	go eventConsumer(pirCh)

	// pirRight.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {

	// 	if p.Get() {

	// 		fmt.Printf("PIR arriving PinRising event")
	// 		// sendEvent(RightRising, sm, &ctx)
	// 		// pirCh <- "RightRising"

	// 	} else {

	// 		fmt.Printf("PIR arriving PinFalling event")
	// 		// sendEvent(RightFalling, sm, &ctx)
	// 		// pirCh <- "RightFalling"
	// 	}

	// })

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
	// green := color.RGBA{0, 255, 0, 255}

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
	for {

		if keyPressed == key0 {
			keyPressed = 0
			log.Printf("key0\n")
		}
		if keyPressed == key1 {
			keyPressed = 0
			log.Printf("key1\n")
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
		// cls(&display)
		msg := fmt.Sprintf("DFT: %v\nAed: %v\nAng: %v\nDed: %v\nDng: %v\nErr: %v\nFls: %v",ctx.DefaultCount,ctx.ArrivedCount,ctx.ArrivedCount,ctx.DepartedCount,ctx.DepartingCount,ctx.ErrorCount,ctx.FalseAlarmCount)
		// tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 10, 30, msg, green)
		time.Sleep(time.Millisecond * 3000)
		fmt.Printf("%v\n",msg)
		fmt.Printf(".")
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

func sendEvent(event fsm.EventID, sm *fsm.StateMachine, ctx *MartyContext) {

	err := sm.SendEvent(event, ctx)
	if err == fsm.ErrEventRejected {
		ctx.ErrorCount += 1
		sm.Current = fsm.Default
	}

}

func eventConsumer(ch chan string) {
	for {
		// Wait for a change in position
		event := <-ch
		fmt.Printf("event: %v ", event)

		if event == "RightRising" {
			sendEvent("RightRising",sm,&ctx)
		}
		if event == "RightFalling" {
			sendEvent("RightFalling",sm,&ctx)
		}		
	}
}
