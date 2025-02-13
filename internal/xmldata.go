package xmldata

import (
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"os"
)

type DefaultStr struct {
	XMLName xml.Name `xml:"Login"`
	Type    string   `xml:"Type"`
}

type LoginAns struct {
	XMLName xml.Name `xml:"Login"`
	Type    string   `xml:"Type"`
	Result  string   `xml:"Result"`
}

type Login struct {
	XMLName xml.Name `xml:"Login"`
	Type    string   `xml:"Type"`
	ID      string   `xml:"ID"`
	PW      string   `xml:"PW"`
}

type LoginResult struct {
	XMLName xml.Name `xml:"Login"`
	Type    string   `xml:"Type"`
	ID      string   `xml:"ID"`
	Result  string   `xml:"Result"`
}

type Chat struct {
	XMLName xml.Name `xml:"Chat"`
	Type    string   `xml:"Type"`
	ID      string   `xml:"ID"`
	Chat    string   `xml:"ChatString"`
}

func XmlInit() {
	testdata := Chat{Type: "Chat", ID: "asdf", Chat: "hansando"}

	file, err := os.Create("data.xml")
	if err != nil {
		fmt.Println("파일 생성 오류", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "   ")
	err = encoder.Encode(testdata)
	if err != nil {
		fmt.Println("XML error: ", err)
		return
	}
	fmt.Println("XML Fin data.xml")
}

func GetDefaultXML(conn net.Conn) (string,[]byte, int){
	fmt.Println("Default XML")

	recv := make([]byte, 4096)

	n, err := conn.Read(recv)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection is closed from client")
			return "", nil, 0
		}
		fmt.Println("False to receive default XML : ", err)
		return "", nil, 0
	}
	if n > 0 {
		fmt.Println("Received raw data: ", string(recv[:n]))
		///////////////
		//xml parsing
		var msg DefaultStr
		err = xml.Unmarshal(recv[:n], &msg)
		if err != nil {
			fmt.Println("Error parsing XmL:", err)
			return "", nil, 0
		}
		fmt.Printf("\nDefault ParsingData\n")
		fmt.Printf("Message Type : %s\n", msg.Type)
		return msg.Type, recv, n
	}
	return "", nil, 0
}

//////////////
