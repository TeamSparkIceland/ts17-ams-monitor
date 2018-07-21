package main

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/tarm/serial"
	"time"
)

func GetPort() string {

	log.Infof("Looking for port")

	for i := 0; i < 256; i++ {

		port := fmt.Sprintf("COM%d", i)

		result := portScanListen(port, BaudRate)
		if result {
			return port
		}
	}

	log.Infof("No port found")
	return ""
}

func portScanListen(port string, baudrate int ) bool {

	s,  err := connectPortScan(port, baudrate)

	if err != nil {
		return false
	} else {

		data := make([]byte, 0, 2048)
		buf := make([]byte, 1)

		for i := 0; i < 50; i++ {

			_, err := s.Read(buf)

			if err != nil {
				log.Errorf("Failed to read serial port data: %v", err)
				s.Close()
				return false
			}

			if buf[0] == '\r' {
				continue
			}

			b := buf[0]
			if b == Delimiter || b == PACKET_START_CHAR || b == STATUS_START_CHAR || b == CURRENT_START_CHAR || b == TSAL_START_CHAR {
				if len(data) > 0 {
					log.Infof("Found data at port %s", port)
					s.Close()
					return true
				}
				data = make([]byte, 0, 2048)
			} else {
				data = append(data, buf[0])
			}
		}
	}
	s.Close()
	return false
}

func connectPortScan(port string, baudrate int) (*serial.Port, error) {
	c := &serial.Config{ Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * 50}
	return serial.OpenPort(c)
}