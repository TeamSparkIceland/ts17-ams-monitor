package main

import (
	"fmt"
	"github.com/golang-ui/nuklear/nk"
	"log"
)

const (
	packFrameWidth  = 100
	packFrameHeight = 300
)

func filler(ctx *nk.Context, height float32) {
	nk.NkLayoutRowStatic(ctx, height, 0, 0)
}

func makeSidebarFrame(ctx *nk.Context, state *State, x, y, w, h float32) {
	windowFlags := (nk.Flags)(nk.WindowNoScrollbar)
	nk.NkPlatformNewFrame()
	bounds := nk.NkRect(x, y, w, h)

	update := nk.NkBegin(ctx, "Sidebar", bounds, windowFlags)

	if update > 0 {
		filler(ctx, 10)

		// Prenta út straum
		nk.NkLayoutRowDynamic(ctx, 14, 1)
		{
			// Búa til label
			nk.NkLabel(ctx, "Current", nk.TextCentered)
			// Búa til lit
			currentColor := nk.NkRgba(102, 178, 255, 255)
			// Búa tit label með texta "x.xx A" með litnum
			nk.NkLabelColored(ctx, fmt.Sprintf("%.2f A", state.Current), nk.TextCentered, currentColor)
		}

		filler(ctx, 10)

		// Segja AMS að senda gögn eða ekki
		nk.NkLayoutRowDynamic(ctx, 60, 1)
		{
			if state.RequestData {
				//nk.NkButtonSymbolLabel(ctx, nk.SymbolRectSolid, "Request data", nk.TextCentered)
				if nk.NkButtonSymbolLabel(ctx, nk.SymbolRectSolid, "Request data", nk.TextCentered) > 0 {
					state.RequestData = !state.RequestData
					log.Println("Requesting data off, ", state.RequestData)
				}
			} else {
				if nk.NkButtonSymbolLabel(ctx, nk.SymbolRectOutline, "Request data", nk.TextCentered) > 0 {
					state.RequestData = !state.RequestData
					log.Println("Requesting data on, ", state.RequestData)
				}
			}
		}

		filler(ctx, 10)

		nk.NkLayoutRowDynamic(ctx, 60, 1)
		{

			//if nk.NkButtonLabel(ctx, "CONN") > 0 {
			//}

			if state.DischargeRequested {
				if nk.NkButtonSymbolLabel(ctx, nk.SymbolRectSolid, "Discharge", nk.TextCentered) > 0 {
					state.DischargeRequested = !state.DischargeRequested
				}
			} else {
				if nk.NkButtonSymbolLabel(ctx, nk.SymbolRectOutline, "Discharge", nk.TextCentered) > 0 {
					state.DischargeRequested = !state.DischargeRequested
				}
			}

		}
		if state.DischargeRequested {
			nk.NkLayoutRowDynamic(ctx, 14, 1)
			{
				nk.NkLabel(ctx, "Target Voltage", nk.TextCentered)
				color := nk.NkRgba(102, 178, 255, 255)
				nk.NkLabelColored(ctx, fmt.Sprintf("%.2f V", state.DischargeTargetVoltage), nk.TextCentered, color)
			}
		}
	}
	nk.NkEnd(ctx)
}

func makeSegmentFrame(ctx *nk.Context, state *State, segmentId int, x, y, w, h float32) {
	windowFlags := (nk.Flags)(nk.WindowBorder | nk.WindowNoScrollbar | nk.WindowTitle)
	frameId := fmt.Sprintf("SEGMENT %d", segmentId)
	nk.NkPlatformNewFrame()
	bounds := nk.NkRect(x, y, w, h)

	update := nk.NkBeginTitled(ctx, frameId, frameId, bounds, windowFlags)
	if update > 0 {

		totalVoltage := state.PackData[(segmentId*2)].Voltage + state.PackData[(segmentId*2)+1].Voltage

		nk.NkLayoutRowDynamic(ctx, 30, 1)
		{
			//		nk.NkStylePushFont(ctx, largeFont.Handle())
			nk.NkLabel(ctx, fmt.Sprintf("%.2f V", totalVoltage), nk.TextCentered)
			//		nk.NkStylePopFont(ctx)
		}

		nk.NkLayoutRowDynamic(ctx, 260, 1)
		{

			for packId := 0; packId < 2; packId++ {
				nk.NkGroupBegin(ctx, fmt.Sprintf("S%d-P%d", segmentId, packId), nk.WindowBorder|nk.WindowNoScrollbar|nk.WindowTitle)
				makePackLayout(ctx, state, (segmentId*2)+packId)
				nk.NkGroupEnd(ctx)
			}

		}
	}
	nk.NkEnd(ctx)
}

func makePackLayout(ctx *nk.Context, state *State, packId int) {
	for i, cell := range state.PackData[packId].Cells {
		nk.NkLayoutRowDynamic(ctx, 14, 3)
		{
			nk.NkLabel(ctx, fmt.Sprintf("C%d", i), nk.TextLeft)
			if cell.DischargeActive {
				color := nk.NkRgba(102, 178, 255, 255)
				nk.NkLabelColored(ctx, fmt.Sprintf("%.2f V", cell.Voltage), nk.TextLeft, color)
			} else {
				nk.NkLabel(ctx, fmt.Sprintf("%.2f V", cell.Voltage), nk.TextLeft)
			}
			nk.NkLabel(ctx, fmt.Sprintf("%.2f C", cell.Temperature), nk.TextLeft)
		}
	}
}

func makePackFrame(ctx *nk.Context, config *Config, state *State, idFont *nk.Font, i int, x, y, w, h float32) {
	windowFlags := (nk.Flags)(nk.WindowBorder | nk.WindowNoScrollbar | nk.WindowTitle)

	frameId := fmt.Sprintf("BMS %d", i)

	nk.NkPlatformNewFrame()

	bounds := nk.NkRect(x, y, w, h)

	if nk.NkBeginTitled(ctx, frameId, frameId, bounds, windowFlags) > 0 {

		makePackLayout(ctx, state, i)
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		{
			nk.NkLabel(ctx, fmt.Sprintf("%.2f V", state.PackData[i].Voltage), nk.TextCentered)
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
