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

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
