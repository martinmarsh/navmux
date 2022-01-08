/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"time"
)

func serial(config *ConfigData) {
	fmt.Println("started navdata serial")
}

func serialreader(port string, name string, baud int64, output []string, channels map[string](chan string)) {
	ever := true
	str := ""
	i := 1
	for ever {
		str = fmt.Sprintf("%s %s %d %d", port, name, baud, i)
		i += 1
		for _, out := range output {
			channels[out] <- str
			time.Sleep(10 * time.Millisecond)
		}
	}
}

/*
	   //	for i, v := range inputs {
	   //		fmt.Println(i, v)
	   		m := viper.GetStringMapString(v)
	   		fmt.Println(m)
	   		var baud int64 = 4800
	   		baud, err := strconv.ParseInt(m["baud"], 10, 64)
	   		if err != nil {
	   			baud = 4800
	   		}
	   		outstr := v + "." + "output"
	   		if m["type"] == "serial" {
	   			fmt.Println("serial", baud, m["name"])
	   		}
	   		output := viper.GetStringSlice(outstr)
	   		fmt.Println(output)

	   		for i2, out := range output {
	   			fmt.Println(i2, out)
	   			if _, ok := channels[out]; !ok {
	   				channels[out] = make(chan string, 10)
	   			}
	   		}
	   		go serial(v, m["name"], baud, output, channels)

	   	}

	   	cont := true

	   	for cont {
	   		prompt := promptui.Select{
	   			Label: "Select Action",
	   			Items: []string{"Continue", "Exit", "Status"},
	   		}

	   		_, result, err := prompt.Run()

	   		if err != nil {
	   			fmt.Printf("Prompt failed %v\n", err)
	   		}

	   		fmt.Printf("You choose %q\n", result)
	   	}
	   	return inputs
	*/
