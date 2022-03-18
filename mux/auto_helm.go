package mux

import (
	"fmt"
	
)

func autoHelmProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	input := config["input"][0]
	go helm(name, input, channels)

}

func helm(name string, input string, channels *map[string](chan string)){
	for {
		str := <-(*channels)[input]
		fmt.Printf("Recieved helm %s\n", str)
	}
}
