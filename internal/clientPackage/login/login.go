package Login

import (
	"fmt"
	"net"
	//Globals "github.com/jenjer/ChatGo/internal/clientPackage"
	xmldata "github.com/jenjer/ChatGo/internal"
	sendxml "github.com/jenjer/ChatGo/internal/clientPackage/xmlSend"
)

func tryLogin(id string, pw string, conn net.Conn) {
	var reqLoginXml xmldata.Login
	reqLoginXml.ID = id
	reqLoginXml.PW = pw
	reqLoginXml.Type = "Login"
	sendxml.SendMessage(conn, reqLoginXml)
}

func Login(conn net.Conn) {
	var ID string
	var PW string
	fmt.Print("Input your Login ID : ")
	fmt.Scan(&ID)

	fmt.Print("Input your Login PW : ")
	fmt.Scan(&PW)

	fmt.Printf("ID : %s, PW : %s\n", ID, PW)

	tryLogin(ID, PW, conn)
	//Globals.SetID(ID)
}
