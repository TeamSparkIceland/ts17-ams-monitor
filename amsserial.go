package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"time"
)

const (
	delimiter = '\n'
)

func connect(port string, baudrate int) *serial.Port {
	c := &serial.Config{Name: port, Baud: baudrate}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Failed to connect to port %s: %v", port, err)
	}
	return s
}

func listen(port string, baudrate int, received chan<- string, outbound <-chan string, quitC <-chan struct{}) {
	s := connect(port, baudrate)
	defer s.Close()

	data := make([]byte, 0, 2048)
	buf := make([]byte, 1)

	for {

		_, err := s.Read(buf)
		if err != nil {
			log.Errorf("Failed to read serial port data: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if buf[0] == '\r' {
			continue
		}

		if buf[0] == delimiter {
			if len(data) > 0 {
				received <- string(data[:])
			}
			data = make([]byte, 0, 2048)
		} else {
			data = append(data, buf[0])
		}

		select {
		case o := <-outbound:
			log.Infof("Sending %s to serial", o)
			s.Write([]byte(o))
		case <-quitC:
			return
		default:
		}
	}
}
