/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package cmd

import (
	"github.com/martinmarsh/nmea-mux"
)

func load_and_run() {
	n := nmea_mux.NewMux()

	// for default config.yaml in current folder
	// optional parameters define folder, filename, format, "config as a string
	if err := n.LoadConfig(); err == nil {

		n.Run()        // Run the virtual devices / go tasks
		n.WaitToStop() // Wait for ever?

	}

}
