package main

import (
	"fmt"
	"github.com/golang-ui/nuklear/nk"
)

const (
	packFrameWidth  = 100
	packFrameHeight = 300
)

func makePackFrame(ctx *nk.Context, config *Config, state *State, idFont *nk.Font, i int, x, y, w, h float32) {
	windowFlags := (nk.Flags)(nk.WindowBorder | nk.WindowNoScrollbar | nk.WindowTitle)

	frameId := fmt.Sprintf("BMS %d", i)

	nk.NkPlatformNewFrame()

	bounds := nk.NkRect(x, y, w, h)

	if nk.NkBeginTitled(ctx, frameId, frameId, bounds, windowFlags) > 0 {

		for i, cell := range state.PackData[i].Cells {
			nk.NkLayoutRowDynamic(ctx, 14, 3)
			{
				nk.NkLabel(ctx, fmt.Sprintf("C%d", i), nk.TextLeft)
				nk.NkLabel(ctx, fmt.Sprintf("%.2f V", cell.Voltage), nk.TextLeft)
				nk.NkLabel(ctx, fmt.Sprintf("%.2f C", cell.Temperature), nk.TextLeft)
			}
		}

	}

	nk.NkEnd(ctx)
}

func makeLogFrame(ctx *nk.Context, state *State, x, y, w, h float32) {
	windowFlags := nk.WindowBorder | nk.WindowTitle

	title := "Serial Log"

	nk.NkPlatformNewFrame()

	bounds := nk.NkRect(x, y, w, h)

	var hide int32
	hide = flag(state.HideDataLog)

	if nk.NkBeginTitled(ctx, title, title, bounds, (nk.Flags)(windowFlags)) > 0 {

		nk.NkLayoutRowDynamic(ctx, 18, 1)
		{
			if nk.NkCheckboxLabel(ctx, "Hide data-packets", &hide) > 0 {
				if hide == 0 {
					state.HideDataLog = false
				} else {
					state.HideDataLog = true
				}
			}
		}

		nk.NkLayoutRowDynamic(ctx, 10, 1)
		{
			for _, line := range state.LogData {
				nk.NkLabel(ctx, line, nk.TextLeft)
			}
		}

	}

	nk.NkEnd(ctx)
}

func makeThresholdViewFrame(ctx *nk.Context, config *Config, state *State, x, y, w float32) {

	windowFlags := nk.WindowBorder | nk.WindowNoScrollbar

	title := "Threshold Settings"
	lineHeight := 15
	linePadding := 6
	lineCount := 4

	nk.NkPlatformNewFrame()

	bounds := nk.NkRect(x, y, 250, float32(lineCount*(lineHeight+linePadding)))

	if nk.NkBeginTitled(ctx, title, title, bounds, (nk.Flags)(windowFlags)) > 0 {

		nk.NkLayoutRowDynamic(ctx, 15, 2)
		{
			nk.NkLabel(ctx, "Voltage (MAX): ", nk.TextLeft)
			nk.NkLabel(ctx, fmt.Sprintf("%0.1f V", config.CellConfig.VoltageMax), nk.TextRight)
		}

		nk.NkLayoutRowDynamic(ctx, 15, 2)
		{
			nk.NkLabel(ctx, "Voltage (MIN): ", nk.TextLeft)
			nk.NkLabel(ctx, fmt.Sprintf("%0.1f V", config.CellConfig.VoltageMin), nk.TextRight)
		}

		nk.NkLayoutRowDynamic(ctx, 15, 2)
		{
			nk.NkLabel(ctx, "Temperature (MAX): ", nk.TextLeft)
			nk.NkLabel(ctx, fmt.Sprintf("%0.1f C", config.CellConfig.TemperatureMax), nk.TextRight)
		}

		nk.NkLayoutRowDynamic(ctx, 15, 2)
		{
			nk.NkLabel(ctx, "Temperature (MIN): ", nk.TextLeft)
			nk.NkLabel(ctx, fmt.Sprintf("%0.1f C", config.CellConfig.TemperatureMin), nk.TextRight)
		}

	}

	nk.NkEnd(ctx)
}
