/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/
package nmea

import (
	"github.com/martinmarsh/nmea0183"
	"fmt"
) 

var Sentences nmea0183.Sentences

func Setup() *nmea0183.Handle{

	err := Sentences.Load()
	if err != nil{
		fmt.Println(fmt.Errorf("**** Error config: %w", err))
		Sentences.SaveLoadDefault()
	}
	handle := Sentences.MakeHandle()
	return handle
}
