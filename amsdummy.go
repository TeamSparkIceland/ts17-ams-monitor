package main

import (
	"fmt"
	"math/rand"
)

func generateDataPackets() []string {
	result := make([]string, 0, 12)
	for packId := 0; packId < 12; packId++ {
		packLine := fmt.Sprintf("D%d|", packId)
		for cellId := 0; cellId < 12; cellId++ {
			voltage := int((rand.Float32() + 3.1) * 1000)
			temperature := int((rand.Float32() * 40) + 15)
			packLine += fmt.Sprintf("%d|%d|%d|", cellId, voltage, temperature)
		}
		result = append(result, packLine)
	}
	return result
}

func dummyListen(received chan string, quitC chan struct{}) {
	for {
		select {
		case <-quitC:
			return
		default:
		}
		for _, line := range generateDataPackets() {
			received <- line
		}
	}
}
