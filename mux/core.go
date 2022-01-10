/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"time"
	"github.com/manifoldco/promptui"
)

type ConfigData struct {
	Index    map[string]([]string)
	TypeList map[string]([]string)
	Values   map[string]map[string]([]string)
}

func Execute(config *ConfigData) {
	channels := make(map[string](chan string))
	fmt.Println("Navmux execute")
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

	for processType, names := range config.TypeList {
		fmt.Println(processType, names)
		switch processType {
			case "serial":
				for _, name := range names {
					serialProcess(name, config.Values[name], &channels)
				}
		}
	}
	for{
		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{"Continue", "Exit", "Status"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}

		fmt.Printf("You choose %q\n", result)
		time.Sleep(5 * time.Second)
		if result == "Exit" {
			break
		}
	}
}
