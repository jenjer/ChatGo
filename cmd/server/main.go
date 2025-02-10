package main

import (
	"fmt"
	"io"
	"net"
)

func handler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())
	recv := make([]byte, 4096)

	for {
		n, err := conn.Read(recv)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection is closed from client : ", conn.RemoteAddr().String())
				return
			}
			fmt.Println("Failed to receive data :", err)
			return
		}
		if n > 0 {
			fmt.Println(string(recv[:n]))
			_, err = conn.Write(recv[:n])
			if err != nil {
				fmt.Printf("Failed to send data to %s: %v\n", conn.RemoteAddr().String(), err)
				return
			}

		}
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil{
		fmt.Println("Fale to Listen : ", err)
	}

	defer l.Close()

	for {
		fmt.Println("Server started on port 8000");
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("failed to Accept %s: ", err)
			continue
		}

		go handler(conn)
	}
}
