package sendXML

import (
	"net"
	xmldata "github.com/jenjer/ChatGo/internal"
	"encoding/xml"
	"fmt"
)

func SendMessage(conn net.Conn, msg xmldata.Login) {
	encoder := xml.NewEncoder(conn)
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Println("err send xml :" , err)
	}
}
