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
	server_port := config["server_port"][0]
	server_ip := config["server_ip"][0]
	input_channel := config["input"][0]
	fmt.Println(server_port, server_ip, input_channel)

	service := server_ip + ":" + server_port

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		fmt.Println("Could not open udp ")
	}

	fmt.Printf("Established connection to %s \n", service)
	fmt.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	fmt.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	// write a message to server
	message := []byte("$GPZDA,110910.59,15,09,2020,00,00*6F\n")

	_, err = conn.Write(message)

	if err != nil {
		fmt.Println(err)
	}

}
