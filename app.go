package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func dataReceiver(state *State, received <-chan string) {
	for {
		select {
		case r := <-received:
			if len(state.LogData) >= LOG_LINES_MAX {
				i := (len(state.LogData) - LOG_LINES_MAX) + 1
				state.LogData = state.LogData[i:]
			}
			state.LogData = append(state.LogData, r)

			if isDataPacket(r) {
				parseDataPacket(r, state)
			}
			if isStatusPacket(r) {
				parseStatusPacket(r, state)
			}
			if isCurrentPacket(r) {
				parseCurrentPacket(r, state)
			}
		}
	}
}

func dataSender(state *State, outbound chan<- string) {
	lastDischargeState := state.DischargeRequested
	for {
		newDischargeState := state.DischargeRequested
		if newDischargeState == lastDischargeState {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if newDischargeState {
			outbound <- "E"
		} else {
			outbound <- "D"
		}
		lastDischargeState = newDischargeState
	}
}

func app() {
	config := NewDefaultConfig()
	state := NewDefaultState()

	quitC := make(chan struct{})
	received := make(chan string)
	outbound := make(chan string)

	if viper.GetBool("dryrun") {
		log.Warn("Using dummy data instead of live data (dryrun option)")
		go dummyListen(received, quitC)
	} else {
		go listen(viper.GetString("port"), 115200, received, outbound, quitC)
	}

	go dataReceiver(state, received)
	go dataSender(state, outbound)

	state.PackData = createPackData(0.0, 0.0, false, false)

	guiLoop(config, state)
}
