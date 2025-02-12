package Login

import (
	"net"
	"fmt"
	"io"
	//xmlstruct "github.com/jenjer/ChatGo/internal"
)

func TryLogin(conn net.Conn)(bool) {
	fmt.Printf("Read ID/PW")

	recv := make([]byte, 4096)

	for {
		n, err := conn.Read(recv)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection is closed from client")
				return false
			}
			fmt.Println("Failed to receive login data : ",err)
			return false
		}

		if n > 0 {
			//print holedata
			fmt.Println("Received raw data: ", string(recv[:n]))
		}
		return true
	}
	return true
}
