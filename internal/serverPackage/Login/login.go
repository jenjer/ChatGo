package Login

import (
	"net"
	"fmt"
	"io"
	"encoding/xml"
	xmlstruct "github.com/jenjer/ChatGo/internal"
	DBConn "github.com/jenjer/ChatGo/internal/serverPackage/DB"
)

func TryLogin(conn net.Conn)(bool) {
	fmt.Printf("Read ID/PW")

	recv := make([]byte, 4096)

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

		//xml parsing
		var msg xmlstruct.Login
		err = xml.Unmarshal(recv[:n], &msg)
		if err != nil {
			fmt.Println("Error parsing XmL:", err)
			return false
		}
		fmt.Printf("\nParsingData\n")
		fmt.Printf("Message Type : %s\n", msg.Type)
		fmt.Printf("Message ID : %s\n", msg.ID)
		fmt.Printf("Message PW : %s\n", msg.PW)

		if (DBConn.ValidateUser(msg.ID, msg.PW) == true) {
			return true
		}
	}
	return false
}
