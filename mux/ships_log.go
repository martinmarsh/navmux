package mux

import (
	"navmux/nmea"
	"fmt"
	"os"
	"time"
	"bufio"
	"encoding/json"
)



func shipsLogProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	f, err := os.Create("ships_log.txt")
	w := bufio.NewWriter(f)

	if err != nil {
		fmt.Println("FATAL Error logging: " + name)
		time.Sleep(time.Minute)
	}

	input := config["input"][0]
	go fileLogger(name, w, input, channels)

}


func fileLogger(name string, writer *bufio.Writer, input string, channels *map[string](chan string)){
	const every = 100
	count := every
	handle := nmea.Setup()
	for {
		str := <-(*channels)[input]
		// fmt.Printf("Recieved log %s\n", str)
		if len(str)>5 && len(str)<89 && str[0] == '$'{
			// fmt.Printf("counter is %d\n", count)
	        count -= 1
			handle.Parse(str)
			if count == 0 {
				count = every
				
				data_map := handle.GetMap()
				data_json, _ := json.Marshal(data_map)
				
				rec_str := fmt.Sprintf("%s\n", string(data_json))
				fmt.Println(rec_str)

				_, err := writer.WriteString(rec_str)
				if err != nil {
					fmt.Println("FATAL Error on write" + name)
				}
				writer.Flush()
			}
		}	
	}
}
