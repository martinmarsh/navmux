package mux

import (
	"fmt"
	"math"
	"navmux/buffer"
	"navmux/io"
	"strconv"
	"strings"

)

func relative_direction(diff float32) float32 {
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
	var p_factor float32 = 100.0
	var i_factor float32 = 2
	var d_factor float32 = 50
	var g_factor float32 = 100
	var base_gain float32 = 0.00001

	input := config["input"][0]
	if p, e := strconv.ParseFloat(config["p_factor"][0], 32); e == nil {
		p_factor = float32(p)
	}
    if i_f, e := strconv.ParseFloat(config["i_factor"][0], 32); e == nil {
		i_factor = float32(i_f)
	}
    if d_f, e := strconv.ParseFloat(config["d_factor"][0], 32); e == nil {
		d_factor = float32(d_f)
	}
    if gain_factor, e := strconv.ParseFloat(config["gain_factor"][0], 32); e == nil{
		g_factor = float32(gain_factor)
	}
	
	go helm(name, input, channels, p_factor,i_factor, d_factor, g_factor, base_gain)
	
}

func helm(
	name string, 
	input string,
	channels *map[string](chan string),
	p_factor float32,
	i_factor float32,
	d_factor float32,
	g_factor float32,
	base_gain float32,
){
	buffer_p := buffer.MakeFloatBuffer(20)
	buffer_d := buffer.MakeFloatBuffer(20)

	var course_to_steer float32
	var heading float32 = 0.0
	var rate_of_course_error_rate float32 = 0.0
	var prev_course_error_rate float32 = 0.0
	var prev_course_error float32 = 0.0
	var auto_on bool = false
	var course_error float32 = 0.0
	var course_error_rate float32 = 0.0
	var gain float32 = g_factor * base_gain

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
				hd, _ := strconv.ParseFloat(parts[1], 32)
				heading = float32(hd)
				//rate of change over 1.5 seconds given constant heading readings of 10ps
				course_error = relative_direction(course_to_steer - heading)
				if buffer_p.Count >= 15 {
					prev_course_error = buffer_p.Read()
					course_error_rate = course_error - prev_course_error
					//fmt.Printf("Heading %.2f, %.2f %.3f\n", hd, prev_head, delta)
				}
				buffer_p.Write(course_error)
				// rate of change of the course error rate
				if buffer_d.Count >= 15 {
					prev_course_error_rate =  buffer_d.Read()
					rate_of_course_error_rate = course_error_rate - prev_course_error_rate
				}
				buffer_d.Write(course_error_rate)
			}
		} else if str == "compute" && auto_on {
			power_calc(course_error, course_error_rate, rate_of_course_error_rate, p_factor, i_factor, d_factor, gain)
		
		} else if len(str) > 2 && str[0] == 'P'{
			if value, e := cmd_value(str); e == nil {
				p_factor = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'D'{
			if value, e := cmd_value(str); e == nil {
				d_factor = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'I'{
			if value, e := cmd_value(str); e == nil {
				i_factor = value
				io.Beep("1l")
			}
		} else if len(str) > 2 && str[0] == 'G'{
			if value, e := cmd_value(str); e == nil {
				g_factor = value
				gain = value * base_gain + 1
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
						io.Helm('X', 0)
						io.Beep("1l")
						fmt.Printf("Auto factor; gain:%.7f  p: %.1f i:%.1f d: %.1f\n", gain, p_factor, i_factor, d_factor)
					}
				}
			case '7':
				if value, e := cmd_value(str); e == nil {
					if value == 0 {
						course_to_steer = heading
						auto_on = true
						power_calc(course_error, course_error_rate, rate_of_course_error_rate, p_factor, i_factor, d_factor, gain)
						io.Beep("3s")
					}
				}
			case '8':
				if value, e := cmd_value(str); e == nil {
					g_factor += value
					gain = g_factor * base_gain
					io.Beep("5s")
				}
			case '2':
				if value, e := cmd_value(str); e == nil {
					if g_factor > value {
						g_factor -= value
						gain = g_factor * base_gain
						io.Beep("4s")
					}
				}
			case '5':
				if value, e := cmd_value(str); e == nil {
					gain = value * base_gain * 100 + 1
					io.Beep("1s")
					io.Beep("1l")
				}
			case '9':
				if value, e := cmd_value(str); e == nil {
					d_factor += value
					io.Beep("2l")
				}
			case '3':
				if value, e := cmd_value(str); e == nil {
					d_factor -= value
					io.Beep("1l")
				}
			}
		}

	}
}

func cmd_value(str string) (float32, error) {
	if p, e := strconv.ParseFloat(str[1:], 32); e == nil {
		return float32(p), nil
	} else {
		return 0, nil
	}
}

func power_calc(course_error, course_error_rate, rate_of_course_error_rate, p_factor, i_factor, d_factor, gain float32){
	/*
	Control is based on PID principles. However, in our case we don't get feedback from the contol surface
	but from the change of direction of the boat.  A constant compass error would produce an output from
	a PID which would continue to drive the rudder position and in effect integrate the turning force.
	To compensate for the in built integrator we differentiate all inputs from the streering compass

	There are 4 response control factors:
	p_factor - proportional gain
	i_factor - integral gain
	d_factor - differential gain
	gain - overall gain used to control sensitivity and scale the PID to values required by the motor driver
	
	*/
	fmt.Printf("course error %.2f, course rate %.3f, raterate %.4f , gain %.6f\n", course_error, course_error_rate, rate_of_course_error_rate, gain)
	// Error correction is a constant value which integrates so is scaled by I
	integration := course_error * i_factor
	// The rate of course_error change is in fact proportional
	proportional := course_error_rate * p_factor
	// The rate of change of the course_error_rate is the differential in our case
	differential := rate_of_course_error_rate * d_factor
	fmt.Printf("integration %.2f, proportional %.3f, diff %.4f \n", integration, proportional, differential)

	power := (integration + proportional + differential) * gain
	power_manager(power)
}

func power_manager(power float32){
	const prefix string = "LR"
	pi := 0			
	if power  > 0 {
		pi = 1
	}

	io.Helm(prefix[pi], math.Abs(float64(power)))
}

