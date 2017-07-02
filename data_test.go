package main

import (
	"testing"
)

func testInitialPackDataCreation(t *testing.T) {
	volt := 1.0
	temp := 2.0
	voltPec := true
	tempPec := false
	packData := createPackData(volt, temp, voltPec, tempPec)

	for p := 0; p < PACKS_IN_SYSTEM; p++ {
		for c := 0; c < CELLS_IN_PACK; c++ {
			if packData[p].Cells[c].Voltage != volt {
				t.Errorf(
					"Voltage in pack %d, cell %d is not correct: %f != %f (expected)",
					p,
					c,
					packData[p].Cells[c].Voltage,
				)
			}
		}
	}
}
