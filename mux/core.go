/*
Copyright © 2024 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"time"
	"os/exec"
)

type ConfigData struct {
	Index    map[string]([]string)
	TypeList map[string]([]string)
	Values   map[string]map[string]([]string)
}

var Monitor_channel chan string
var Udp_monitor_active bool

func Monitor(str string, print bool, udp bool) {
	if udp && Udp_monitor_active {
		Monitor_channel <- str
	}
	if print {
		fmt.Println(str)
	}
}

func Execute(config *ConfigData) {
	// wait for everything to connect on boot up
	time.Sleep(5 * time.Second)
	
	channels := make(map[string](chan string))
	fmt.Println("Navmux execute")
	channels["command"] = make(chan string, 2)
	
	for name, param := range config.Index {
		for _, value := range param {
			if value == "outputs" {
				for _, chanName := range config.Values[name][value] {
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 30)
					}
				}
			}
			if value == "input" {
				for _, chanName := range config.Values[name][value] {
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 30)
					}
				}
			}
		}
	}

	for processType, names := range config.TypeList {
		fmt.Println(processType, names)
		for _, name := range names {
			switch processType {
			case "serial":
				serialProcess(name, config.Values[name], &channels)
			case "udp_client":
				udpClientProcess(name, config.Values[name], &channels)
			case "processor":
				processorProcess(name, config.Values[name], &channels)
			case "udp_listen":
				udpListenerProcess(name, config.Values[name], &channels)
			case "compass_resolver":
				compassResolveProcess(name, config.Values[name], &channels)
 	
			}
		}
	}

	
	for {
		command := <-(channels["command"])
		fmt.Printf("Command '%s' received\n", command)
		
	}
}


func Exit() {
	out, err := exec.Command("shutdown","-h","now").Output()

    // if there is an error with our execution
    // handle it here
    if err != nil {
        fmt.Printf("%s", err)
    }
    // as the out variable defined above is of type []byte we need to convert
    // this to a string or else we will see garbage printed out in our console
    // this is how we convert it to a string
    fmt.Println("Command Successfully Executed")
    output := string(out[:])
    fmt.Println(output)

}