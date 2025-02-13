package Login

import (
	"net"
	"fmt"
//	"io"
	"encoding/xml"
	xmlstruct "github.com/jenjer/ChatGo/internal"
	DBConn "github.com/jenjer/ChatGo/internal/serverPackage/DB"
)

func TryLogin(conn net.Conn, DbConn *DBConn.UserDB)(bool,string) {
	fmt.Printf("Read ID/PW")

	
	xmltype, recv, n := xmlstruct.GetDefaultXML(conn)
	if xmltype != "Login" || recv == nil {
		fmt.Println("type is not Login")
		return false, "";
	} 

	//print holedata
	fmt.Println("Received raw data: ", string(recv[:n]))

	//xml parsing
	var msg xmlstruct.Login
	err := xml.Unmarshal(recv[:n], &msg)
	if err != nil {
		fmt.Println("Error parsing XmL:", err)
		return false, ""
	}
	fmt.Printf("\nLogin ParsingData\n")
	fmt.Printf("Message Type : %s\n", msg.Type)
	fmt.Printf("Message ID : %s\n", msg.ID)
	fmt.Printf("Message PW : %s\n", msg.PW)

	if temp, err := DbConn.ValidateUser(msg.ID, msg.PW); temp == true {
		return true, msg.ID
	} else {
		fmt.Println("Something is wrong : " ,err)
	}
	return false, ""
}
