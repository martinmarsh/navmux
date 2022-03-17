/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package io

import (
	"fmt"
	"time"
	"strconv"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

var Beep_channel chan string

func Init() error {
	Beep_channel = make(chan string, 2)

	go beeperTask()

	fmt.Println("Loading periph.io drivers")
	// Load periph.io drivers:
	if _, err := host.Init(); err != nil {
		return err
	}
	return nil
}

func Beep(style string){
	Beep_channel <- style
}

func beeperTask(){
	for{
		beep := <- Beep_channel
		if len(beep) == 2 {
			count, _ := strconv.Atoi(string(beep[0]))
			
			length := 400
			if beep[1] == 's'{
				length = 100
			} else if beep[1] == 'l'{
				length = 800
			}
			for l := 0; l < count; l++  {
				t := time.NewTicker(time.Duration(length) * time.Millisecond)
				rpi.P1_22.Out(gpio.High)
				<-t.C
				rpi.P1_22.Out(gpio.Low)
				<-t.C
			}
		}
	}

}
