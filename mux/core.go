/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"navmux/io"
	"os"
	"time"
	"os/exec"

	"github.com/stianeikeland/go-rpio/v4"
)

type ConfigData struct {
	Index    map[string]([]string)
	TypeList map[string]([]string)
	Values   map[string]map[string]([]string)
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

	if err := rpio.Open(); err != nil {
		fmt.Println("RPIO - could not be openned")
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	 if err := io.Init(&channels); err != nil {
		fmt.Println("Failed to set up gpio")
		panic("error GPIO")
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

	
	io.Beep("1s")
	
	for {
		command := <-(channels["command"])
		fmt.Printf("Command '%s' received\n", command)
		if command == "99" {
			io.Beep("2l")
			Exit()
		} else if len(command) > 1 && len(command) <=5 {
			switch command[0]{
			case '1':
				channels["to_helm"] <- "P" + command[1:]
			case '2':
				channels["to_helm"] <- "I" + command[1:]
			case '3':
				channels["to_helm"] <- "D" + command[1:]
			case '4':
				channels["to_helm"] <- "G" + command[1:]
			}
			
				
		}
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