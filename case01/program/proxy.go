package main

import (
	"fmt"
	"io"
	"net"
)

func proxyForward(from io.Reader, to io.Writer) error{

	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	if(fromIsWriter && toIsReader){
		go func(){
			io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)
	return err
}

func main() {
	listener, err := net.Listen("tcp", "localhost:9999")

	if err != nil {
		return
	}

	for{
		conn, err := listener.Accept()

		if err != nil {
			return
		}

		go func(from net.Conn){
			fmt.Println("Accepted")

			to, err := net.Dial("tcp", "localhost:1111")

			if err != nil {
				return
			}

			err = proxyForward(from, to)

			if err != nil {
				return
			}
			
		}(conn)
	}
}