package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	PACKET_START_CHAR     = 'D'
	PACKET_DELIMITER      = "|"
	VALID_DELIMITER_COUNT = 49
)

func isDataPacket(packet string) bool {
	if packet[0] != PACKET_START_CHAR {
		return false
	}

	if strings.Count(packet, PACKET_DELIMITER) != VALID_DELIMITER_COUNT {
		log.Warnf("Datapacket has invalid delimiter count: %s\n", packet)
		return false
	}

	return true
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
	for i := 0; i < 12; i++ {
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
