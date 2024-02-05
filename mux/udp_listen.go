/*
Copyright © 2024 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"net"

)

func udpListenerProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	// listens on a port and writes to output channels
	
	server_port := config["port"][0]
	to_chans := ""
	for _, out := range config["outputs"] {
		to_chans += fmt.Sprintf(" %s,", out)
	}
	Monitor(fmt.Sprintf("Upd_listen; name: %s  Port: %s channels: %s",name, server_port, to_chans), true, true)

	tag:= ""

	if config["origin_tag"] != nil {
		tag = fmt.Sprintf("@%s@", config["origin_tag"])
	}

	if len(config["outputs"]) > 0 {
		go udpListener(name, server_port, config["outputs"], tag,  channels)
	}
	
}

func udpListener(name string, server_port string, outputs []string, tag string, channels *map[string](chan string)) {
	const maxBufferSize = 1024
	pc, err := net.ListenPacket("udp", "0.0.0.0:"+server_port)
	if err != nil {
		Monitor(fmt.Sprintf("Error; Upd_listen; action: ABORTED, error: %s", err.Error()), true, true)
		return
	}
	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	
	for{
		n, _, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("packet error")
			Monitor(fmt.Sprintf("Error; Upd_listen; Packet Error; action: ignored, error: %s", err.Error()), true, true)
			return
	
		} else {
			for _, out := range outputs {
				(*channels)[out] <- tag + string(buffer[:n])
			}
		}

	}


}
