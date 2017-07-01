package main

import (
	"github.com/golang-ui/nuklear/nk"
)

type CellConfig struct {
	VoltageMax     float32
	VoltageMin     float32
	TemperatureMax float32
	TemperatureMin float32
}

type Config struct {
	CellConfig          CellConfig
	CellBackgroundColor nk.Color
}

type Cell struct {
	Voltage             float32
	Temperature         float32
	VoltagePecError     bool
	TemperaturePecError bool
	DischargeActive     bool
}

type Pack struct {
	Cells   []Cell
	Voltage float32
}

type State struct {
	HideDataLog bool
	LogData     []string
	PackData    []Pack
}

func NewDefaultConfig() *Config {
	return &Config{
		CellConfig: CellConfig{3.0, 4.2, 5.0, 55.0},
	}
}
