package main

import (
	"fmt"
	"net"
)

func main() {
	dial, err := net.Dial("tcp", "localhost:9999")

	if err != nil {
		return
	}

	defer dial.Close()

	payload := Binary("Quack Quack!")

	_, err = payload.WriteTo(dial)

	if err != nil {
		return
	}

	p, err := Decode(dial)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Connection timeout!")
		} else {
			fmt.Println(err)
		}
		return 
	}

	fmt.Println(p)
}