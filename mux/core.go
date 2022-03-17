/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	
)

type ConfigData struct {
	Index    map[string]([]string)
	TypeList map[string]([]string)
	Values   map[string]map[string]([]string)
}

func Execute(config *ConfigData) {
	channels := make(map[string](chan string))
	fmt.Println("Navmux execute")
	channels["command"] = make(chan string, 2)
	channels["to_beeper"] = make(chan string, 2)
	for name, param := range config.Index {
		for _, value := range param {
			if value == "outputs" {
				for _, chanName := range config.Values[name][value] {
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 10)
					}
				}
			}
		}
	}

	if err := initGPIO(); err != nil {
		fmt.Println("gpio error")
	}

	for processType, names := range config.TypeList {
		fmt.Println(processType, names)
		for _, name := range names {
			switch processType {
			case "serial":
				serialProcess(name, config.Values[name], &channels)
			case "udp_client":
				udpClientProcess(name, config.Values[name], &channels)
			case "keyboard":
				keyBoardProcess(name, config.Values[name], &channels)
			case "ships_log":
				shipsLogProcess(name, config.Values[name], &channels)
			case "auto-helm":
				autoHelmProcess(name, config.Values[name], &channels)
					
			}
		}
	}

	go beeperTask(&channels)


	channels["to_udp_client"] <- "$GPZDA,110910.59,15,09,2020,00,00*6F"
	
	channels["to_udp_client"] <- "$HCHDM,172.5,M*28"
	channels["to_udp_client"] <- "$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"
	channels["to_udp_client"] <- "$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C"
	channels["to_udp_client"] <- "$SSDPT,2.8,-0.7"
	channels["to_log"] <- "$GPZDA,110910.59,15,09,2020,00,00*6F"
	channels["to_log"] <- "$HCHDM,172.5,M*28"
	channels["to_log"] <- "$GPZDA,110910.59,15,09,2020,00,00*6F"
	channels["to_log"] <- "$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C"
	channels["to_log"] <- "$SSDPT,2.8,-0.7"

	channels["to_beeper"] <- "1s"

	for {
		command := <-(channels["command"])
		fmt.Printf("Command '%s' received\n", command)

	}


}
