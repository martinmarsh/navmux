package mux

import (
	"bufio"
	"encoding/json"
	"fmt"
	"navmux/nmea"
	"os"
	"time"
	"strings"

	"github.com/martinmarsh/nmea0183"
)



func shipsLogProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	input := config["input"][0]
	go fileLogger(name, input, channels)

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

	if len(str)>5 && len(str)<89 && str[0] == '$'{
		// fmt.Printf("counter is %d\n", count)
		_, _, error := handle.Parse(str)
		return error
	}
	return fmt.Errorf("%s", "no leading dollar")
}
