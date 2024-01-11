/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package io

import (
	"fmt"
	"math"
	"strconv"
	"time"

	//"github.com/stianeikeland/go-rpio/v4"
	//"periph.io/x/periph/conn/gpio"
	//"periph.io/x/periph/host"
	//"periph.io/x/periph/host/rpi"
	//"periph.io/x/conn/v3/physic"
)

type helm_control struct {
    control byte
    power  float64
}

var Beep_channel chan string
var Motor_channel chan helm_control

var beep_pin = rpio.Pin(25)
var left_pin = rpio.Pin(24)
var right_pin = rpio.Pin(23)
var power_pin = rpio.Pin(18)


func Init(channels *map[string](chan string)) error {
	Beep_channel = make(chan string, 4)
	Motor_channel = make(chan helm_control, 3)

	go beeperTask()
	go helmTask(channels)

	//fmt.Println("Loading periph.io drivers")
	// Load periph.io drivers:
	//if _, err := host.Init(); err != nil {
	//	return err
	//}
	//if err := rpio.Open(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//defer rpio.Close()

	beep_pin.Output()
	left_pin.Output()
	right_pin.Output()
	power_pin.Output()

	left_pin.Low()
	right_pin.Low()
	power_pin.Low()   

	return nil
}

func Beep(style string){
	Beep_channel <- style
}

func Helm(control byte, power float64){
	// control = 1 if power is a signed value
	// control = 0 to turn off motor
	// control = 'L' or 'R' if power is unsigned
	if control == 1 {
		control = 'R'
		if power < 0 {
			control = 'L'
			power = math.Abs(power)
		}

	}
	message :=  helm_control{
		control: control,
		power: power,   	// <=1 1 = 100%
	}
	fmt.Println(message)
	Motor_channel <- message
}

func helmTask(channels *map[string](chan string)){
	const max_power = 3000    // 3ms cycle time  300us min
	const max_loops = 85	 // 85 x 3  255 ms read channel period
	const max_power_slow = 20000    // 20ms cycle time  3ms min
	const max_loops_slow = 14	 // 14 x 20  280 ms read channel period
	t1 := time.Duration(0)
	t2 := time.Duration(max_power_slow) * time.Microsecond
	mp := max_power_slow
	ml := max_loops_slow
		
	for {	
		//ask helm to compute another value
		(*channels)["to_helm"] <-"compute"

		select {
		case motor := <- Motor_channel:
			fmt.Println("motor control")
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
			case 1:
			default:
				left_pin.Low()
				right_pin.Low()
				//rpi.P1_18.Out(gpio.Low)
				//rpi.P1_16.Out(gpio.Low)
			}
			p1 := 0
			mp = max_power_slow
			ml = max_loops_slow

			if motor.power > 0.2 && motor.power < 0.8 {	
				mp = max_power
				ml = max_loops
			}
		    if motor.power > 0.95 {
				p1 = mp
			} else if motor.power < 0.02{
				p1 = 0
			} else {
				p1 = int(float64(mp) * motor.power)
			}
			t1 = time.Duration(p1) * time.Microsecond
			t2 = time.Duration(mp - p1) * time.Microsecond
			//fmt.Printf("%d %d\n", t1, t2)	
		default:
			// continue
		}

		if t1 == 0 {
			power_pin.Low()
			time.Sleep(250 * time.Millisecond)
		} else if t2 == 0 {
			power_pin.High()
			time.Sleep(250 * time.Millisecond)
		} else {
			for i := 0; i < ml; i++ {
				if t1 > 0 {
					power_pin.High()
					time.Sleep(t1)
				}
				if t2 > 0 {
					power_pin.Low()
					time.Sleep(t2)
				}
			}
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
