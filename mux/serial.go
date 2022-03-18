/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

func serialProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	fmt.Println("started navmux serial " + name)
	baud, err := strconv.ParseInt(config["baud"][0], 10, 64)
	if err != nil {
		baud = 4800
	}
	mode := &serial.Mode{
		BaudRate: int(baud),
	}
	portName := config["name"][0]
	port, err := serial.Open(portName, mode)
	if err != nil {
		fmt.Println("no serial port " + portName)
	} else {
		if len(config["outputs"]) > 0 {
			fmt.Println("Open read serial port " + portName)
			go serialReader(name, port, config["outputs"], channels)
			
		}
		if len(config["input"]) > 0 {
			fmt.Println("Open write serial port " + portName)
			go serialWriter(name, port, config["input"], channels)
		}

	}

}

func serialReader(name string, port serial.Port, outputs []string, channels *map[string](chan string)) {
	buff := make([]byte, 100)
	time.Sleep(3 * time.Second)
	for {
		n, err := port.Read(buff)
		
		if err != nil {
			fmt.Println("FATAL Error on port " + name)
			time.Sleep(time.Minute)
		}
		if n == 0 {
			fmt.Println("\nEOF on read of " + name)
			time.Sleep(time.Minute)
		}
		
		// If we receive a newline send to output channels
		if strings.Contains(string(buff[:n]), "\n") {
			str := string(buff[:n])
			for _, out := range outputs {
				(*channels)[out] <- str
			}
		}
	}

}

func serialWriter(name string, port serial.Port, input []string, channels *map[string](chan string)) {
	time.Sleep(2 * time.Second)
	for {
		for _, in := range input {
			str := <-(*channels)[in]
			_, err := port.Write([]byte(str))
			if err != nil {
				fmt.Println("FATAL Error on port" + name)
				time.Sleep(time.Minute)
			}
			
		}
	}

}
