package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateDataPackets() []string {
	result := make([]string, 0, 10)
	for packId := 0; packId < 10; packId++ {
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

func generateStatusPacket() string {
	voltage := int((rand.Float32() + 3.1) * 1000)
	return fmt.Sprintf("S|%d|", voltage)
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
	}
}
