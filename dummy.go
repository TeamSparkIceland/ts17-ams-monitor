package main

import (
	"time"
)

func dummyListen(received chan string, quitC chan struct{}) {
	for {
		select {
		case <-quitC:
			return
		case <-time.After(time.Second * 1):
			received <- "STRING\n"
		}
	}
}
