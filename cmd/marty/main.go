package main

import (
	"fmt"
	"github.com/tonygilkerson/marty/pkg/fsm"
	"log"
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

func main() {
	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("no impl here yet, run this instead:\ngo test -v ./...\n")

}

type MartyContext struct {
	DefaultCount    int
	ArrivedCount    int
	ArrivingCount    int
	DepartedCount   int
	DepartingCount   int
	ErrorCount      int
	FalseAlarmCount int
}

func (c *MartyContext) String() string {
  cCopy := *c
	return fmt.Sprintf("MartyContext: %+v\n", cCopy)
}

//
// DefaultAction
//
type DefaultAction struct{}

func (a *DefaultAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DefaultCount += 1

	log.Printf("DefaultAction\n\n")
	return fsm.NoOp
}

//
// ArrivedAction
//
type ArrivedAction struct{}

func (a *ArrivedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ArrivedCount += 1

	log.Printf("ArrivedAction\n")
	return Reset
}

//
// ArrivingAction
//
type ArrivingAction struct{}

func (a *ArrivingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ArrivingCount += 1

	log.Printf("ArrivingAction\n")
	return fsm.NoOp
}

//
// DepartedAction
//
type DepartedAction struct{}

func (a *DepartedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DepartedCount += 1

	log.Printf("DepartedAction\n")
	return Reset
}

//
// DepartingAction
//
type DepartingAction struct{}

func (a *DepartingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.DepartingCount += 1

	log.Printf("DepartingAction\n")
	return fsm.NoOp
}

//
// ErrorAction
//
type ErrorAction struct{}

func (a *ErrorAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MartyContext)
	ctx.ErrorCount += 1

	log.Printf("ErrorAction\n")
	return fsm.NoOp
}

//
// FalseAlarmAction
//
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
					LeftRising: Departing,
				},
			},

			Arriving: fsm.State{
				Action: &ArrivingAction{},
				Events: fsm.Events{
					LeftRising: Arrived,
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


