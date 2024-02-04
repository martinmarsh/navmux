package mux

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"

)




func logProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	var generators = make(map[string] *GENERATE)
	input := config["input"][0]
	go fileLogger(name, input, channels)
	fmt.Println(config)
	for dotKey, value := range config {
		key := strings.Split(dotKey, ".")
		if key[0] == "0183_generators" {
			if generators[key[1]] == nil {
				generators[key[1]] = &GENERATE{sentence: key[1]}
				generators[key[1]].alternatives = make(map[string] *ALTERNATIVE)
			}
			for j := 2; j < len(key); j++{
				switch key[j]{
				case "every":
					generators[key[1]].every, _ = strconv.Atoi(value[0])
				case "prefix":
					generators[key[1]].prefix = value[0]
				case "send_to":
					generators[key[1]].send_to = value
				case "alternative":
					if generators[key[1]].alternatives[key[j+1]] == nil {
						generators[key[1]].alternatives[key[j+1]] = &ALTERNATIVE{variable: key[j+1]}
					}
				case "replace_with":
					if generators[key[1]].alternatives[key[j-1]] == nil {
						generators[key[1]].alternatives[key[j-1]] = &ALTERNATIVE{variable: key[j+1]}
					}
					generators[key[1]].alternatives[key[j-1]].replace_with = value[0]
				case "if":
					fmt.Printf("if %s - %s\n", key[j-1], value) 
				case "and":
					fmt.Printf("and %s - %s\n", key[j-1], value)
				case "or":
					fmt.Printf("or %s - %s\n", key[j-1], value)

				default:
					fmt.Printf("missed %s - %s\n", key[j], key[1]) 
				}
			
			}
		}

	}
	for n, g:= range generators{
		fmt.Printf("gen %s every %d, prefix %s sentence %s goto %s \n", n, g.every, g.prefix, g.sentence, g.send_to)
		for a, ga := range g.alternatives{
			fmt.Printf("alternatives %s  %s replace with %s\n", a, ga.variable, ga.replace_with)
		}


	}
}


func fileLogger(name string, input string, channels *map[string](chan string)){
	var writer *bufio.Writer
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	handle := nmea.Setup()
	file_closed := true
	
	for {
		select {
		case str := <-(*channels)[input]:
			parse(str, handle)
		case <-ticker.C:
			data_map := handle.GetMap()
			
				
			if file_closed {
				if dt, ok := data_map["datetime"]; ok {
					dt = strings.Replace(dt[:16], ":", "_", -1)
					file_name := fmt.Sprintf("ships_log_%s.txt", dt)
					f, err := os.Create(file_name)
					writer = bufio.NewWriter(f)
					if err != nil {
						fmt.Println("FATAL Error logging: " + name)
						time.Sleep(time.Minute)
					} else {
						file_closed = false
					}

				}

			} else {
				data_json, _ := json.Marshal(data_map)
				rec_str := fmt.Sprintf("%s\n", string(data_json))
				//fmt.Println(rec_str)
				_, err := writer.WriteString(rec_str)
				if err != nil {
					fmt.Println("FATAL Error on write" + name)	
					writer.Flush()
				}
				writer.Flush()
			}
		}	
	}
}

func parse(str string, handle *nmea0183.Handle) error{

	defer func(){
		if r := recover(); r != nil {
			str = ""
			fmt.Println("\n** Recover from NEMEA Panic **")
		}
	}()
	str = strings.TrimSpace(str)
	if len(str)>5 && len(str)<89 && str[0] == '$'{
		// fmt.Printf("counter is %d\n", count)
		_, _, error := handle.Parse(str)
		return error
	}
	return fmt.Errorf("%s", "no leading dollar")
}
