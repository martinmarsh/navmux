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

	RemoteAddr, _ := net.ResolveUDPAddr("udp", server_addr)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	//defer conn.Close()

	if err != nil {
		fmt.Printf("Could not open udp server %s/n", name)
	} else {
		fmt.Printf("Established connection to %s \n", server_addr)
		fmt.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
		fmt.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

		go udpWriter(name, conn, input_channel, channels)

	}

}

func udpWriter(name string, conn *net.UDPConn, input string, channels *map[string](chan string)) {
	for {
		str := <-(*channels)[input]
		_, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Println("FATAL Error on UDP connection" + name)
		}
		
	}
}
