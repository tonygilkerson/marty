package main

// To run tests 
// $ go test -v ./...
//

import (
	"github.com/tonygilkerson/marty/pkg/fsm"
	"testing"
)

func TestMartyStateMachine(t *testing.T) {

	// Create a new instance of the light switch state machine.
	martyFSM := newMartyFSM()
	var ctx MartyContext

	//
	// A car arriving
	//
	ctx = resetContext()
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(LeftRising, martyFSM, &ctx)

	if ctx.DefaultCount == 1 &&
		ctx.ArrivedCount == 1 &&
		ctx.ArrivingCount == 1 &&
		ctx.DepartedCount == 0 &&
		ctx.DepartingCount == 0 &&
		ctx.ErrorCount == 0 &&
		ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car arriving\nexpected: {DefaultCount:1 ArrivedCount:1 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", ctx)
	}

	//
	// A car departing
	//
	ctx = resetContext()
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(RightRising, martyFSM, &ctx)

	if ctx.DefaultCount == 1 &&
		ctx.ArrivedCount == 0 &&
		ctx.ArrivingCount == 0 &&
		ctx.DepartedCount == 1 &&
		ctx.DepartingCount == 1 &&
		ctx.ErrorCount == 0 &&
		ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car departing\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:0 DepartedCount:1 DepartingCount:1 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", ctx)
	}


	//
	// FalseAlarm from the Arriving direction
	// A car approaching but stops short, turns around, backups up or something
	//
	ctx = resetContext()
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(RightFalling, martyFSM, &ctx)

	if ctx.DefaultCount == 1 &&
		ctx.ArrivedCount == 0 &&
		ctx.ArrivingCount == 1 &&
		ctx.DepartedCount == 0 &&
		ctx.DepartingCount == 0 &&
		ctx.ErrorCount == 0 &&
		ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Arriving direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", ctx)
	}

	//
	// FalseAlarm from the Departing direction
	//
	ctx = resetContext()
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(LeftFalling, martyFSM, &ctx)

	if ctx.DefaultCount == 1 &&
		ctx.ArrivedCount == 0 &&
		ctx.ArrivingCount == 0 &&
		ctx.DepartedCount == 0 &&
		ctx.DepartingCount == 1 &&
		ctx.ErrorCount == 0 &&
		ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Departing direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", ctx)
	}


	//
	// Error from the Departing direction
	// Error - should never get two Rising events in a row from the same direction
	//
	ctx = resetContext()
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(LeftRising, martyFSM, &ctx)

	if ctx.DefaultCount == 0 &&
		ctx.ArrivedCount == 0 &&
		ctx.ArrivingCount == 0 &&
		ctx.DepartedCount == 0 &&
		ctx.DepartingCount == 1 &&
		ctx.ErrorCount == 1 &&
		ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Departing direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:0 DepartedCount:0 DepartingCount:1 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", ctx)
	}

	//
	// Error from the Arriving direction
	// Error - should never get two Rising events in a row from the same direction
	//
	ctx = resetContext()
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(RightRising, martyFSM, &ctx)

	if ctx.DefaultCount == 0 &&
		ctx.ArrivedCount == 0 &&
		ctx.ArrivingCount == 1 &&
		ctx.DepartedCount == 0 &&
		ctx.DepartingCount == 0 &&
		ctx.ErrorCount == 1 &&
		ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Arriving direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", ctx)
	}

	//
	// Combination of events
	//
	ctx = resetContext()
	
	// A car arriving
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(LeftRising, martyFSM, &ctx)

	// A car departing
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(RightRising, martyFSM, &ctx)

	// FalseAlarm from the Arriving direction
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(RightFalling, martyFSM, &ctx)

	// FalseAlarm from the Departing direction
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(LeftFalling, martyFSM, &ctx)

	// Error from the Departing direction
	sendEvent(LeftRising, martyFSM, &ctx)
	sendEvent(LeftRising, martyFSM, &ctx)

	// Error from the Arriving direction
	sendEvent(RightRising, martyFSM, &ctx)
	sendEvent(RightRising, martyFSM, &ctx)

	if ctx.DefaultCount == 4 &&
		ctx.ArrivedCount == 1 &&
		ctx.ArrivingCount == 3 &&
		ctx.DepartedCount == 1 &&
		ctx.DepartingCount == 3 &&
		ctx.ErrorCount == 2 &&
		ctx.FalseAlarmCount == 2 {
		// all good
	} else {
		t.Errorf("A combination of events\nexpected: {DefaultCount:4 ArrivedCount:1 ArrivingCount:3 DepartedCount:1 DepartingCount:3 ErrorCount:2 FalseAlarmCount:2}\ngot:      %+v", ctx)
	}

}

func resetContext() MartyContext {
	return MartyContext{
		DefaultCount:    0,
		ArrivedCount:    0,
		ArrivingCount:   0,
		DepartedCount:   0,
		DepartingCount:  0,
		ErrorCount:      0,
		FalseAlarmCount: 0,
	}
}

func sendEvent(event fsm.EventID, sm *fsm.StateMachine, ctx *MartyContext) {

	err := sm.SendEvent(event, ctx)
	if err == fsm.ErrEventRejected {
		ctx.ErrorCount += 1
		sm.Current = fsm.Default
	}

}


