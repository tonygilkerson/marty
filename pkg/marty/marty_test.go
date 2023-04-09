package marty

// To run tests
// $ go test -v ./...
//

import (
	"testing"
)

func TestMartyStateMachine(t *testing.T) {

	// Create a new instance of the light switch state machine.
	m := New()

	//
	// A car arriving
	//
	m.ResetContext()
	m.SendEvent(ArriveRising)
	m.SendEvent(DepartRising)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 1 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car arriving\nexpected: {DefaultCount:1 ArrivedCount:1 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// A car departing
	//
	m.ResetContext()
	m.SendEvent(DepartRising)
	m.SendEvent(ArriveRising)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 1 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car departing\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:0 DepartedCount:1 DepartingCount:1 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// FalseAlarm from the Arriving direction
	// A car approaching but stops short, turns around, backups up or something
	//
	m.ResetContext()
	m.SendEvent(ArriveRising)
	m.SendEvent(ArriveFalling)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Arriving direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", m.Ctx)
	}

	//
	// FalseAlarm from the Departing direction
	//
	m.ResetContext()
	m.SendEvent(DepartRising)
	m.SendEvent(DepartFalling)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Departing direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", m.Ctx)
	}

	//
	// Error from the Departing direction
	// Error - should never get two Rising events in a row from the same direction
	//
	m.ResetContext()
	m.SendEvent(DepartRising)
	m.SendEvent(DepartRising)

	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Departing direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:0 DepartedCount:0 DepartingCount:1 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// Error from the Arriving direction
	// Error - should never get two Rising events in a row from the same direction
	//
	m.ResetContext()
	m.SendEvent(ArriveRising)
	m.SendEvent(ArriveRising)

	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Arriving direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// Default goes to Default if LD or RD
	//
	m.ResetContext()
	m.SendEvent(DepartRising)
	m.SendEvent(ArriveRising)
	m.SendEvent(ArriveFalling)
	m.SendEvent(DepartFalling)

	if m.Ctx.DefaultCount == 3 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 1 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error default goes to default\nexpected: {DefaultCount:3 ArrivedCount:0 ArrivingCount:0 DepartedCount:1 DepartingCount:1 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// Combination of events
	//
	m.ResetContext()

	// A car arriving
	m.SendEvent(ArriveRising)
	m.SendEvent(DepartRising)

	// A car departing
	m.SendEvent(DepartRising)
	m.SendEvent(ArriveRising)

	// FalseAlarm from the Arriving direction
	m.SendEvent(ArriveRising)
	m.SendEvent(ArriveFalling)

	// FalseAlarm from the Departing direction
	m.SendEvent(DepartRising)
	m.SendEvent(DepartFalling)

	// Error from the Departing direction
	m.SendEvent(DepartRising)
	m.SendEvent(DepartRising)

	// Error from the Arriving direction
	m.SendEvent(ArriveRising)
	m.SendEvent(ArriveRising)

	if m.Ctx.DefaultCount == 4 &&
		m.Ctx.ArrivedCount == 1 &&
		m.Ctx.ArrivingCount == 3 &&
		m.Ctx.DepartedCount == 1 &&
		m.Ctx.DepartingCount == 3 &&
		m.Ctx.ErrorCount == 2 &&
		m.Ctx.FalseAlarmCount == 2 {
		// all good
	} else {
		t.Errorf("A combination of events\nexpected: {DefaultCount:4 ArrivedCount:1 ArrivingCount:3 DepartedCount:1 DepartingCount:3 ErrorCount:2 FalseAlarmCount:2}\ngot:      %+v", m.Ctx)
	}

}
