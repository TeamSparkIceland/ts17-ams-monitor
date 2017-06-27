package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
)

var RootCmd = &cobra.Command{
	Use:   "ts17-ams",
	Short: "TS17 Accumulator Management System Monitor",
	Run: func(cmd *cobra.Command, args []string) {
		app()
	},
}

func init() {
	runtime.LockOSThread()
	RootCmd.Flags().StringP("port", "p", "COM1", "Serial port")
	RootCmd.Flags().BoolP("dryrun", "d", false, "Fake data instead of serial communication")
	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("dryrun", RootCmd.Flags().Lookup("dryrun"))
	viper.SetDefault("port", "COM1")
	viper.SetDefault("dryrun", false)
}

func createPackData() []Pack {
	result := make([]Pack, 12)
	for i := 0; i < 12; i++ {
		cells := make([]Cell, 12)
		for j := 0; j < 12; j++ {
			cells[i] = Cell{
				Voltage:             3.95,
				Temperature:         25.0,
				VoltagePecError:     false,
				TemperaturePecError: false,
			}
		}
		result[i] = Pack{
			Cells: cells,
		}
	}
	return result
}

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
		}
	}
}

func main() {

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func app() {
	config := NewDefaultConfig()
	state := &State{
		HideDataLog: false,
		LogData:     make([]string, 0, LOG_LINES_MAX),
	}

	quitC := make(chan struct{})
	received := make(chan string)

	if viper.GetBool("dryrun") {
		log.Warn("Using dummy data instead of live data (dryrun option)")
		go dummyListen(received, quitC)
	} else {
		go listen(viper.GetString("port"), 115200, received, quitC)
	}

	go dataReceiver(state, received)

	state.PackData = createPackData()

	guiLoop(config, state)
}
