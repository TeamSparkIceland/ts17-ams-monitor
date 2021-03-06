package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateDataPackets() []string {
	result := make([]string, 0, 12)
	for packId := 0; packId < 12; packId++ {
		packLine := fmt.Sprintf("D%d|", packId)
		for cellId := 0; cellId < 12; cellId++ {
			voltage := int((rand.Float32() + 3.1) * 1000)
			temperature := int((rand.Float32() * 40) + 15)
			discharge := 0
			packLine += fmt.Sprintf("%d|%d|%d|%d|", cellId, voltage, temperature, discharge)
		}
		result = append(result, packLine)
	}
	return result
}

func generateCurrentPacket() string {
	current := int((rand.Float32() * 400 - 200) * 1000)
	return fmt.Sprintf("C|%d|", current)
}

func generateStatusPacket() string {
	voltage := int((rand.Float32() + 3.1) * 1000)
	return fmt.Sprintf("S|%d|", voltage)
}

func generateTSALPackets() []string {
	result := make([]string, 0, 2)
	for tsalID :=  0; tsalID < 4; tsalID++ {
		statusBitBool := rand.Int() % 2 == 0 // random bool
		statusBit := 0
		if statusBitBool == true {
			statusBit = 1
			}
		tsalLine := fmt.Sprintf("T%d|%d|", tsalID, statusBit)
		result  = append(result, tsalLine)
	}
	return result
}

func dummyListen(received chan string, quitC chan struct{}) {
	for {
		select {
		case <-quitC:
			return
		case <-time.After(1 * time.Second):
		}
		for _, line := range generateDataPackets() {
			received <- line
		}
		received <- generateStatusPacket()
		received <- generateCurrentPacket()

		for _, line := range generateTSALPackets() {
			received <- line
		}
	}
}
