/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"net"

)

func udpClientProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	fmt.Println("started navmux udp " + name)
	server_addr := config["server_address"][0]
	input_channel := config["input"][0]
	fmt.Println(server_addr, input_channel)

	go udpWriter(name, server_addr, input_channel, channels)

}

func udpWriter(name string, server_addr string, input string, channels *map[string](chan string)) {
	connection := false
	for{
		RemoteAddr, _ := net.ResolveUDPAddr("udp", server_addr)
		conn, err := net.DialUDP("udp", nil, RemoteAddr)
		
		if err != nil {
			fmt.Printf("Could not open udp server %s\n", name)
			//ensure channel is cleared then retry
			connection = false
			for i :=0 ; i > 100; i++ {
				<-(*channels)[input]
			}
		} else {
			defer conn.Close()
			fmt.Printf("Established connection to %s \n", server_addr)
			fmt.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
			fmt.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())
			connection = true
		}
		if connection {
			for {
				str := <-(*channels)[input]
				_, err := conn.Write([]byte(str))
				if err != nil {
					fmt.Println("FATAL Error on UDP connection" + name)
				}
				
			}
		}
	}
}
