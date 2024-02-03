/*
Copyright © 2024 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"encoding/json"
	"navmux/nmea"
	"strconv"
)

func compassResolveProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	fmt.Println("started navmux Compass resolve " + name)
	input_channel := config["input"][0]
	fmt.Println(input_channel)

	go compassResolveWriter(name, input_channel, channels)

}

func compassResolveWriter(name string, input string, channels *map[string](chan string)) {
	handle := nmea.Setup()
	for{
		str := <-(*channels)[input]
		fmt.Println("Message got: " + str)
		parse(str, handle)
	
		//handle.Parse("$HCHDM,172.5,M*28")
		data_map := handle.GetMap()
		l := len(data_map["hdm"]) 
		head1 := 0.0
		valid1 := false

		head2, err2 := strconv.ParseFloat(data_map["hdm2"], 64)
		if l >0 {
			h, err1 := strconv.ParseFloat(data_map["hdm"][:l-2], 64)
			if err1 == nil{
				valid1 = true
				head1 = h
			}
		} 
	
		if err2 == nil && data_map["compass_status"] == "3333" {
			data_map["hdm"] = data_map["hdm2"] + "°M"

		}
		fmt.Printf("%s, %s, %.2f, %.2f, %b\n", data_map["hdm2"], data_map["hdm"], head2, head1, valid1)
		data_json, _ := json.Marshal(data_map)
		fmt.Printf("%s\n", string(data_json))
	}
}
