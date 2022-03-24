package mux

import (
	"fmt"
	"navmux/buffer"
	"navmux/io"
	"navmux/pid"
	"strconv"
	"strings"

)

func relative_direction(diff float64) float64 {
    if diff < -180.0 {
        diff += 360.0
	} else if diff > 180.0 {
        diff -= 360.0
	}
    return diff
}

func checksum(s string) string {
	check_sum := 0
	nmea_data := []byte(s)
	for i := 1; i < len(s); i++ {
		check_sum ^= (int)(nmea_data[i])
	}
	return fmt.Sprintf("%02X", check_sum)
}

func autoHelmProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	

	pid := pid.MakePid(1, 0.2, 0.5, 0.00001)

	pid.Scale_gain =100
	pid.Scale_kd = 100
	pid.Scale_ki = 100
	pid.Scale_kp = 100

	input := config["input"][0]
	if p, e := strconv.ParseFloat(config["p_factor"][0], 32); e == nil {
		pid.Scale_kp = p
	}
    if i_f, e := strconv.ParseFloat(config["i_factor"][0], 32); e == nil {
		pid.Scale_ki = i_f
	}
    if d_f, e := strconv.ParseFloat(config["d_factor"][0], 32); e == nil {
		pid.Scale_kd = d_f
	}
    if gain_factor, e := strconv.ParseFloat(config["gain_factor"][0], 32); e == nil{
		pid.Scale_gain = gain_factor
	}

	
	go helm(name, input, channels, pid)
	
}

// Helm takes collects instructions and a compass bearing at 10hz from the input channel
// A PID is used to calculate the actuating signal which is sent to the motor controller.
// The helm motor runs with a speed defined by actuating signal (AS) value either left or right
// rotation; turning the wheel continuously.  The rudder is not driven to a position by the AS
// as would be required if the error signal was based on the course error. The boat would not
// be steering straight ahead at zero error but the rudder would be only just be comming
// to a halt at the maximum deflection.  In short the integration effect would make steering 
// unstable at any proportional gain. To overcome this issue  we consider that the rudder position
// can be sensed by the rate of turn of the boat. Zero rate of turn is straight ahead.
//
// We therefore use the rate of course change as the feedback signal to calculate the error for the
// PID.  It would not make sense to subtract the rate of turn from the desired heading so the desired
// heading is converted to a desired rate of turn.  This has 2 benfits: 
// 1) The turn rate is controlled and is not too excessive for large corrections or tacking
// 2) The rudder effectiveness which varies greatly with boat speed is automatically compensated.
// 
// The PID calculates the AS signal for every Compass input at a constant 10Hz.  The motor calls
// back through a channel when it is ready to receive the next AS instruction.
//
func helm(name string,  input string, channels *map[string](chan string), pid *pid.Pid) {

	buffer_3 := buffer.MakeFloatBuffer(40)
	buffer_1_5 := buffer.MakeFloatBuffer(20)
	buffer_0_5 := buffer.MakeFloatBuffer(10)

	var course_to_steer float64
	var desired_rate_of_turn float64
	var turn_rate float64

	var auto_on bool = false
	var heading float64 = 0.0

	for {
		str := <-(*channels)[input]
		var err error
		fmt.Printf("Received helm command %s\n", str)
		if len(str)> 9 && str[0:6] == "$HCHDM"{
			end_byte := len(str)
			if str[end_byte-3] == '*' {
				check_code := checksum(str[:end_byte-3])
				end_byte -= 2
				if check_code != str[end_byte:] {
					err_mess := fmt.Sprintf("error: %s != %s", check_code, str[end_byte:])
					err = fmt.Errorf("check sum error: %s", err_mess)
				}
				end_byte--
			}
		
			if err == nil{
				parts := strings.Split(str[1:end_byte], ",")
				heading, _ = strconv.ParseFloat(parts[1], 64)
				average_count := 0
				previous := 0.0
		
				if buffer_3.Count >= 30 {
					previous := buffer_3.Read()
					turn_rate = (heading - previous)/3
					average_count++ 
					//fmt.Printf("Heading %.2f, %.2f %.3f\n", hd, prev_head, delta)
				}
				buffer_3.Write(heading)

				if buffer_1_5.Count >= 15 {
					previous = buffer_1_5.Read()
					turn_rate += relative_direction(heading - previous)/1.5 
					average_count++ 
					//fmt.Printf("Heading %.2f, %.2f %.3f\n", hd, prev_head, delta)
				}
				buffer_3.Write(heading)

				if buffer_0_5.Count >= 5 {
					previous = buffer_1_5.Read()
					turn_rate += relative_direction(heading - previous)/0.5 
					average_count++ 
					//fmt.Printf("Heading %.2f, %.2f %.3f\n", hd, prev_head, delta)
				}
				buffer_3.Write(heading)

				if average_count > 0 {
					turn_rate /= float64(average_count)
					desired_rate_of_turn = relative_direction(course_to_steer - heading)/5  //5s to correct
					av := pid.Compute(desired_rate_of_turn - turn_rate )
					io.Helm(1, av)
				}
				
			}
		} else if str == "compute" && auto_on {
			
		
		} else if len(str) > 2 && str[0] == 'P'{
			if value, e := cmd_value(str); e == nil {
				pid.Scale_kp = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'D'{
			if value, e := cmd_value(str); e == nil {
				pid.Scale_kd = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'I'{
			if value, e := cmd_value(str); e == nil {
				pid.Scale_ki = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'G'{
			if value, e := cmd_value(str); e == nil {
				pid.Scale_gain -= value
				io.Beep("1l")
			}
	
		} else if len(str) > 0 && len(str) < 4 {
			// do key commands
			switch str[0] { 
			case '4':
				if value, e := cmd_value(str); e == nil {
					course_to_steer =  relative_direction(course_to_steer - value)
					io.Beep("2s")
				}
			case '6':
				if value, e := cmd_value(str); e == nil {
					course_to_steer =  relative_direction(course_to_steer + value)
					io.Beep("1s")
				}
			case '1':
				if value, e := cmd_value(str); e == nil {
					if value == 0 {
						auto_on = false
						io.Helm(0, 0)
						io.Beep("1l")
						fmt.Printf("Auto factor; gain:%.7f  p: %.1f i:%.1f d: %.1f\n", pid.Scale_gain, pid.Scale_kp, pid.Scale_ki, pid.Scale_kd)
					}
				}
			case '7':
				if value, e := cmd_value(str); e == nil {
					if value == 0 {
						course_to_steer = heading
						auto_on = true
						io.Helm(1, 0)
						io.Beep("3s")
					}
				}
			case '8':
				if value, e := cmd_value(str); e == nil {
					pid.Scale_gain *= value
					io.Beep("5s")
				}
			case '2':
				if value, e := cmd_value(str); e == nil {
					if value > 1 {
						pid.Scale_gain /= value
						io.Beep("4s")
					}
				}
			case '5':
				if value, e := cmd_value(str); e == nil {
					if value == 0 {
						pid.Scale_gain =100
						pid.Scale_kd = 100
						pid.Scale_ki = 100
						pid.Scale_kp = 100
						io.Beep("1s")
						io.Beep("1l")
					}
				}
			case '9':
				if value, e := cmd_value(str); e == nil {
					pid.Scale_ki *= value
					io.Beep("2l")
				}
			case '3':
				if value, e := cmd_value(str); e == nil {
					if value > 1 {
						pid.Scale_ki /= value
						io.Beep("1l")
					}
				}
			}
		}

	}
}

func cmd_value(str string) (float64, error) {
	if p, e := strconv.ParseFloat(str[1:], 64); e == nil {
		return p, nil
	} else {
		return 0, nil
	}
}
