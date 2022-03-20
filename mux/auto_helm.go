package mux

import (
	"fmt"
	"navmux/nmea"
	"strconv"
	
)

func autoHelmProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	input := config["input"][0]
	go helm(name, input, channels)

}

func helm(name string, input string, channels *map[string](chan string)){

	for {
		str := <-(*channels)[input]
		heading := 0
		var prev_head [10]string
		
		fmt.Printf("Received helm command %s\n", str)
		if str[0] == '$'{
			data,  _, sentenceType, err := nmea.Handle.ParseToMap(str)
			if err != nil && sentenceType == "hdm" {
				ch := data["hdm"]
				l := len(ch)
				if l > 3 {
				    heading = strconv.ParseFloat(data[:4], 64)
				}
			}
		}
	}
}
