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
			time.Sleep(time.Second)
		}
		if len(config["input"]) > 0 {
			fmt.Println("Open write serial port " + portName)
			go serialWriter(name, port, config["input"], channels)
			time.Sleep(time.Second)
		}

	}

}

func serialReader(name string, port serial.Port, outputs []string, channels *map[string](chan string)) {
	buff := make([]byte, 100)

	for {
		n, err := port.Read(buff)
		fmt.Println("Data read")
		if err != nil {
			fmt.Println("FATAL Error on port " + name)
			time.Sleep(time.Minute)
		}
		if n == 0 {
			fmt.Println("\nEOF on read of " + name)
			time.Sleep(time.Minute)
		}
		fmt.Printf("%v", string(buff[:n]))
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
	for {
		for _, in := range input {
			str := <-(*channels)[in]
			fmt.Println("Channel input to send via " + name + "Data: " + str)

			n, err := port.Write([]byte(str))
			if err != nil {
				fmt.Println("FATAL Error on port" + name)
				time.Sleep(time.Minute)
			}
			fmt.Printf("Sent %v bytes\n", n)

		}
	}

}
