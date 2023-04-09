package marty

import (
	"fmt"
	"log"

	"github.com/tonygilkerson/marty/pkg/fsm"
)

const (
	// States
	Arriving   fsm.StateID = "Arriving"
	Arrived    fsm.StateID = "Arrived"
	Departing  fsm.StateID = "Departing"
	Departed   fsm.StateID = "Departed"
	FalseAlarm fsm.StateID = "FalseAlarm"
	Error      fsm.StateID = "Error"

	//Events
	ArriveRising  fsm.EventID = "ArriveRising"
	ArriveFalling fsm.EventID = "ArriveFalling"
	DepartRising  fsm.EventID = "DepartRising"
	DepartFalling fsm.EventID = "DepartFalling"
	Reset         fsm.EventID = "Reset"
)

type Context struct {
	DefaultCount    int
	ArrivedCount    int
	ArrivingCount   int
	DepartedCount   int
	DepartingCount  int
	ErrorCount      int
	FalseAlarmCount int
}

type Marty struct {
	StateMachine fsm.StateMachine
	Ctx          Context
}

func (c *Context) String() string {
	cCopy := *c
	return fmt.Sprintf("Context: %+v\n", cCopy)
}

// sendEvent sends an event to the state machine.
func (m *Marty) SendEvent(event fsm.EventID) {

	err := m.StateMachine.SendEvent(event, &m.Ctx)
	if err == fsm.ErrEventRejected {
		m.Ctx.ErrorCount += 1
		m.StateMachine.Current = fsm.Default
	}

}

func (m *Marty) ResetContext() {
	m.Ctx = Context{
		DefaultCount:    0,
		ArrivedCount:    0,
		ArrivingCount:   0,
		DepartedCount:   0,
		DepartingCount:  0,
		ErrorCount:      0,
		FalseAlarmCount: 0,
	}
}

// DefaultAction
type DefaultAction struct{}

func (a *DefaultAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DefaultCount += 1

	log.Printf("DefaultAction\n\n")
	return fsm.NoOp
}

// ArrivedAction
type ArrivedAction struct{}

func (a *ArrivedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ArrivedCount += 1

	log.Printf("ArrivedAction\n")
	return Reset
}

// ArrivingAction
type ArrivingAction struct{}

func (a *ArrivingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ArrivingCount += 1

	log.Printf("ArrivingAction\n")
	return fsm.NoOp
}

// DepartedAction
type DepartedAction struct{}

func (a *DepartedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DepartedCount += 1

	log.Printf("DepartedAction\n")
	return Reset
}

// DepartingAction
type DepartingAction struct{}

func (a *DepartingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DepartingCount += 1

	log.Printf("DepartingAction\n")
	return fsm.NoOp
}

// ErrorAction
type ErrorAction struct{}

func (a *ErrorAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ErrorCount += 1

	log.Printf("ErrorAction\n")
	return fsm.NoOp
}

// FalseAlarmAction
type FalseAlarmAction struct{}

func (a *FalseAlarmAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.FalseAlarmCount += 1

	log.Printf("FalseAlarmAction\n")
	return Reset
}

func New() Marty {

	var marty Marty
	marty.StateMachine = fsm.StateMachine{
		Current:  fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{

			fsm.Default: fsm.State{
				Action: &DefaultAction{},
				Events: fsm.Events{
					ArriveRising:  Arriving,
					DepartRising:  Departing,
					ArriveFalling: fsm.Default,
					DepartFalling: fsm.Default,
				},
			},

			Arriving: fsm.State{
				Action: &ArrivingAction{},
				Events: fsm.Events{
					DepartRising:  Arrived,
					ArriveFalling: FalseAlarm,
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
					DepartFalling: FalseAlarm,
					ArriveRising:  Departed,
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

	return marty
}
