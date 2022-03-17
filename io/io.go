/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package io

import (
	"fmt"
	"time"
	"strconv"
	"os"
	"github.com/stianeikeland/go-rpio/v4"

	//"periph.io/x/periph/conn/gpio"
	//"periph.io/x/periph/host"
	//"periph.io/x/periph/host/rpi"
	//"periph.io/x/conn/v3/physic"
)

type helm_control struct {
    control byte
    pulse  int
	idle int
}

var Beep_channel chan string
var Motor_channel chan helm_control

var beep_pin = rpio.Pin(25)
var left_pin = rpio.Pin(24)
var right_pin = rpio.Pin(23)
var power_pin = rpio.Pin(18)


func Init() error {
	Beep_channel = make(chan string, 4)
	Motor_channel = make(chan helm_control, 3)

	go beeperTask()
	go helmTask()

	//fmt.Println("Loading periph.io drivers")
	// Load periph.io drivers:
	//if _, err := host.Init(); err != nil {
	//	return err
	//}
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	beep_pin.Output()
	left_pin.Output()
	right_pin.Output()

	return nil
}

func Beep(style string){
	Beep_channel <- style
}

func Helm(control byte, on_period int, base_power int){
	// control = R L or X for Left Right or OFF
	message :=  helm_control{
		control: control,
		pulse: on_period,
		idle: base_power,
	}
	left_pin.Low()
	right_pin.Low()
	
    power_pin.Mode(rpio.Pwm)
	power_pin.Freq(64000)
	power_pin.DutyCycle(0, 32)
	//rpi.P1_18.Out(gpio.Low)
	//rpi.P1_16.Out(gpio.Low)
	//rpi.P1_12.Out(gpio.Low)

	Motor_channel <- message
}

func helmTask(){
	for {
		motor := <- Motor_channel
		switch motor.control {
			case 'L':
				right_pin.Low()
				left_pin.High()
				//rpi.P1_16.Out(gpio.Low)
				//rpi.P1_18.Out(gpio.High)
			case 'R':
				left_pin.Low()
				right_pin.High()
				//rpi.P1_18.Out(gpio.Low)
				//rpi.P1_16.Out(gpio.High)
			case 'X':
			default:
				left_pin.Low()
				right_pin.Low()
				//rpi.P1_18.Out(gpio.Low)
				//rpi.P1_16.Out(gpio.Low)
		}
		t := time.NewTicker(2 * time.Second)
		for i := 0; i > 10; i++{
			//rpi.P1_12.PWM(gpio.DutyMax/1, 1000)
			power_pin.DutyCycle(1, 32)
			<-t.C
			//rpi.P1_12.PWM(gpio.DutyMax/9, 1000)
			power_pin.DutyCycle(30, 32)
			<-t.C
		}

	}

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
				//rpi.P1_22.Out(gpio.High)
				beep_pin.High() 
				<-t.C
				//rpi.P1_22.Out(gpio.Low)
				beep_pin.Low() 
				<-t.C
			}
		}
	}

}
