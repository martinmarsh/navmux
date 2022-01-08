/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package cmd

import (
	"fmt"
	"github.com/martinmarsh/navmux/mux"
	"strings"

	"github.com/spf13/viper"
)

/*

func serial(port string, name string, baud int64, output []string, channels map[string](chan string)) {
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
*/

func loadConfig() *mux.ConfigData {

	fmt.Println("\nLoading config")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := /*  */ err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(fmt.Errorf("fatal error config file err: %w not found = %t", err, ok))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error in config file: %w", err))
		}
	}

	//channels := make(map[string](chan string))

	all := viper.AllKeys()

	confData := &mux.ConfigData{
		Index:    make(map[string]([]string)),
		TypeList: make(map[string]([]string)),
		Values:   make(map[string]([]string)),
	}
	//index := make(map[string]([]string))
	//type_list := make(map[string]([]string))

	// Find keys in yaml assumes >2 deep collects map by 1st part of key
	// also find names by device type every section has a type
	for _, k := range all {
		fmt.Println(k)
		key := strings.SplitN(k, ".", 2)
		if _, ok := confData.Values[k]; !ok {
			confData.Values[k] = viper.GetStringSlice(k)
		}

		if key[1] == "type" {
			type_value := viper.GetString(k)
			if _, ok := confData.TypeList[type_value]; !ok {
				confData.TypeList[type_value] = []string{key[0]}
			} else {
				confData.TypeList[type_value] = append(confData.TypeList[type_value], key[0])
			}

		}
		if _, ok := confData.Index[key[0]]; !ok {
			confData.Index[key[0]] = []string{key[1]}
		} else {
			confData.Index[key[0]] = append(confData.Index[key[0]], key[1])
		}
	}

	// Now process find types
	fmt.Println("end")
	return confData

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
}
