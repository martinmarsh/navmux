/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"time"
	"strconv"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

func initGPIO() error {
	fmt.Println("Loading periph.io drivers")
	// Load periph.io drivers:
	if _, err := host.Init(); err != nil {
		return err
	}
	return nil
}

func beeperTask(channels *map[string](chan string)){
	for{
		beep := <-(*channels)["to_beeper"]
		if len(beep) == 2 {
			count, _ := strconv.Atoi(string(beep[0]))
			
			length := 4
			if beep[1] == 's'{
				length = 1
			} else if beep[1] == 'l'{
				length = 8
			}
			for l := 0; l < count; l++  {
				t := time.NewTicker(time.Duration(length) * time.Second /10)
				rpi.P1_22.Out(gpio.High)
				<-t.C
				rpi.P1_22.Out(gpio.Low)
				<-t.C
			}
		}
	}

}
