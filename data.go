package main

// Serial data structures and handling

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	PACKET_START_CHAR       = 'D'
	PACKET_DELIMITER        = "|"
	VALID_DELIMITER_COUNT   = 49 // 1 + 4 * packs = 1 + 4*12 = 1 + 48 = 49 .... 1 + 4 * 10 = 41
	CELLS_IN_PACK           = 12
	PACKS_IN_SYSTEM         = 12
	STATUS_DELIMITER_COUNT  = 2
	STATUS_START_CHAR       = 'S'
	CURRENT_DELIMITER_COUNT = 2
	CURRENT_START_CHAR      = 'C'
	TSAL_DELIMITER_COUNT 	= 2
	TSAL_START_CHAR			= 'T'
)

func validatePacket(packet string, startChar byte, delimiterCount int) bool {
	if packet[0] != startChar {
		return false
	}

	if strings.Count(packet, PACKET_DELIMITER) != delimiterCount {
		log.Warnf("Datapacket has invalid delimiter count: %s\n", packet)
		return false
	}

	return true
}

func isDataPacket(packet string) bool {
	return validatePacket(packet, PACKET_START_CHAR, VALID_DELIMITER_COUNT)
}

func isStatusPacket(packet string) bool {
	return validatePacket(packet, STATUS_START_CHAR, STATUS_DELIMITER_COUNT)
}

func isCurrentPacket(packet string) bool {
	return validatePacket(packet, CURRENT_START_CHAR, CURRENT_DELIMITER_COUNT)
}

func isTSALPacket(packet string) bool {
	return validatePacket(packet, TSAL_START_CHAR, TSAL_DELIMITER_COUNT)
}

func parseTSALPacket(packet string, state *State){
	parts := strings.Split(packet[1:], PACKET_DELIMITER)[1:]
	tsalId, err := strconv.Atoi(packet[1:2])
	if err != nil {
		log.Warnf("Failed to convert TSAL ID %s to integer\n", packet[1:2])
	} else {
		tsalStatus, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Warnf("Failed to convert Tsal Status bit to integer\n", parts[0])
			return
		}

		switch tsalId {
		case 0:
			state.TsalAirNeg = tsalStatus == 1
		case 1:
			state.TsalAirPos = tsalStatus == 1
		case 2:
			state.TsalConnector = tsalStatus == 1
		case 3:
			state.TsalMC = tsalStatus == 1
		}
	}

}

func parseCurrentPacket(packet string, state *State) {
	parts := strings.Split(packet[1:], PACKET_DELIMITER)[1:]
	current, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Warnf("Failed to convert current %s to integer", parts[0])
		return
	}
	state.Current = float32(current) / 1000.0
}

func parseStatusPacket(packet string, state *State) {
	parts := strings.Split(packet[1:], PACKET_DELIMITER)[1:]
	dischargeVoltage, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Warnf("Failed to convert discharge voltage %s to integer", parts[0])
		return
	}
	state.DischargeTargetVoltage = float32(dischargeVoltage) / 1000.0
}

/*
 * Data Packet
 * D<BMS_ID>|<CELL_ID>|<VOLT>|<TEMP>|<DISCHARGE>|<CELL_ID>...|
 */

func parseDataPacket(packet string, state *State) {

	var totalVoltage float32 = 0

	parts := strings.Split(packet[1:], PACKET_DELIMITER)

	// Get BMS index
	bmsId, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Warnf("Invalid BMS id %s , expected an integer\n", packet[0])
		return
	}

	// Traverse Cell indexes
	for i := 0; i < CELLS_IN_PACK; i++ {
		cellId, err := strconv.Atoi(parts[(1 + (4 * i))])
		if err != nil {
			log.Warnf("Failed to convert cell-id %s to integer\n", parts[(1+(3*i))])
			return
		}

		voltagePart := parts[(2 + (4 * i))]
		if voltagePart == "PEC" {
			// Configure pec error
			state.PackData[bmsId].Cells[cellId].Voltage = 0
			state.PackData[bmsId].Cells[cellId].VoltagePecError = true
		} else {
			voltage, err := strconv.Atoi(voltagePart)
			if err != nil {
				log.Warnf("Failed to convert voltage %s to integer\n", voltagePart)
				return
			}
			state.PackData[bmsId].Cells[cellId].Voltage = float32(voltage) / 1000.0
			state.PackData[bmsId].Cells[cellId].VoltagePecError = false
			totalVoltage = totalVoltage + float32(voltage)/1000.0
		}

		temperaturePart := parts[(3 + (4 * i))]
		if temperaturePart == "PEC" {
			// Configure pec error
			state.PackData[bmsId].Cells[cellId].Temperature = 0
			state.PackData[bmsId].Cells[cellId].TemperaturePecError = true
		} else {
			temperature, err := strconv.Atoi(temperaturePart)
			if err != nil {
				log.Warnf("Failed to convert temperature %s to integer\n", temperaturePart)
				return
			}
			state.PackData[bmsId].Cells[cellId].Temperature = float32(temperature)
			state.PackData[bmsId].Cells[cellId].TemperaturePecError = false
		}

		dischargePart := parts[(4 + (4 * i))]
		state.PackData[bmsId].Cells[cellId].DischargeActive = dischargePart == "1"

	}
	state.PackData[bmsId].Voltage = totalVoltage
}

func createPackData(volt, temp float64, voltPec, tempPec bool) []Pack {
	result := make([]Pack, PACKS_IN_SYSTEM)
	for i := 0; i < len(result); i++ {
		cells := make([]Cell, CELLS_IN_PACK)
		for j := 0; j < len(result); j++ {
			cells[j] = Cell{
				Voltage:             float32(volt),
				Temperature:         float32(temp),
				VoltagePecError:     voltPec,
				TemperaturePecError: tempPec,
			}
		}
		result[i] = Pack{
			Cells: cells,
		}
	}
	return result
}
