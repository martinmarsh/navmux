/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
)

type ConfigData struct {
	Index map[string]([]string)
	TypeList map[string]([]string)
	Values map[string]([]string)
}

func Execute(config *ConfigData){
	fmt.Println("started navdata execute")
	serial(config)
}

