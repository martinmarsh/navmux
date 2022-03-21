package mux

import (
	"fmt"
	"navmux/buffer"
	"navmux/io"
	"strconv"
	"math"
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
	input := config["input"][0]
	power_channel := make(chan float32, 1)
	go helm(name, input, channels, &power_channel)
	go power_manager(&power_channel)

}

func helm(name string, input string, channels *map[string](chan string), power_chan *chan float32){
	buffer := buffer.MakeFloatBuffer(12)
	var course_to_steer float32
	var heading float32 = 0.0
	var turn_speed_factor float32 = 100
	var gain float32 = 0.01
	
	for {
		str := <-(*channels)[input]
		var delta float32
		var err error
		//fmt.Printf("Received helm command %s\n", str)
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
				if buffer.Count >= 10 {
					prev_head := buffer.Read()
					delta = heading - prev_head
					//fmt.Printf("Heading %.2f, %.2f %.3f\n", hd, prev_head, delta)
				}
				buffer.Write(heading)
				
			}
		}


		error_correct := relative_direction(course_to_steer - heading)
        turn_rate := relative_direction(delta)

		power := (error_correct - turn_rate * turn_speed_factor) * gain
	   
		*power_chan <- power

		//pi := 0
		//if power > 0 {
		//	pi = 1
		//}
        //fmt.Printf("power %c %f\n",prefix[pi], math.Abs(power))


		//io.Helm(prefix[pi], math.Abs(power))

	}
}


func power_manager(power_chan *chan float32){
	const prefix string = "LR"
    var av_power float32 = 0
	i := 0

	for{
		
		power := <- (*power_chan)
		av_power += power
		i++
		if i > 7 {
			pi := 0
			av_power = av_power / float32(i)
			if av_power  > 0 {
				pi = 1
			}
			//fmt.Printf("power %c %f\n",prefix[pi], math.Abs(power))
			io.Helm(prefix[pi], math.Abs(float64(av_power )))
			av_power = 0
			i = 0
		}
	}

}

