/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
)

type ConfigData struct {
	Index map[string]([]string)
	TypeList map[string]([]string)
	Values map[string]([]string)
}


func Execute(config *ConfigData){
	channels := make(map[string](chan string))
	fmt.Println("started navdata execute")
	for name, param := range config.Index {
		for _, value := range param {
			if value  == "outputs" {
				for _, chanName := range config.Values[name + "." + value]{
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 10)
					}
				}
			}
		}
	}

	for processType, names := range config.TypeList {
		fmt.Println(processType, names)
		switch processType {
			case "serial":
				for _, name := range names {
				  serial(name, config, &channels)
				}

		} 
	}

}

