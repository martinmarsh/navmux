/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/
package main

import (
	"navmux/cmd"
	"navmux/nmea"

) 

func main() {
	nmea.Setup()
	cmd.Execute()
}
