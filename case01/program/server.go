package main

import (
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:1111")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	for{
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		go func(c net.Conn){
			defer func(){
				c.Close()
			}()

			payload, err := Decode(c)

			if err != nil {
				return
			}

			fmt.Println(payload)
			
			p := Binary("Payload " + string(payload.Bytes()) + " Success")

			_, err = p.WriteTo(c)

			if err != nil {
				return
			}

		}(conn)
	}

}