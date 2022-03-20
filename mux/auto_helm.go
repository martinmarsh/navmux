package mux

import (
	"fmt"
	"navmux/nmea"
	"navmux/buffer"
	"navmux/io"
	"strconv"
	"math"
)

func relative_direction(diff float32) float32 {
    if diff < -180.0 {
        diff += 360.0
	} else if diff > 180.0 {
        diff -= 360.0
	}
    return diff
}

func autoHelmProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	input := config["input"][0]
	go helm(name, input, channels)

}

func helm(name string, input string, channels *map[string](chan string)){
	buffer := buffer.MakeFloatBuffer(12)
	var course_to_steer float32
	var heading float32 = 0.0
	var turn_speed_factor float32 = 100
	var gain float32 = 0.01
	const prefix string = "LR"

	for {
		str := <-(*channels)[input]
		var delta float32
		fmt.Printf("Received helm command %s\n", str)
		if str[0] == '$'{
			data,  _, sentenceType, err := nmea.Handle.ParseToMap(str)
			if err != nil && sentenceType == "hdm" {
				ch := data["hdm"]
				l := len(ch)
				if l > 3 {
				    hd, _ := strconv.ParseFloat(ch[:4], 32)
					heading = float32(hd)
					if buffer.Count >= 10 {
						prev_head := buffer.Read()
						delta = heading - prev_head
					}
					buffer.Write(heading)
				}
			}
		}


		error_correct := relative_direction(course_to_steer - heading)
        turn_rate := relative_direction(delta)

		power := float64((error_correct - turn_rate * turn_speed_factor) * gain)


		pi := 0
		if power > 0 {
			pi = 1
		}

		io.Helm(prefix[pi], math.Abs(power))

	}
}

