package xmldata

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Login struct {
	XMLName xml.Name	`xml:"Login"`
	Type	string		`xml:"Type"`
	ID		string		`xml:"ID"`
	PW		string		`xml:"PW"`
}


type Chat struct {
	XMLName xml.Name	`xml:"Chat"`
	Type	string		`xml:"Type"`
	ID		string		`xml:"ID"`
	Chat	string		`xml:"ChatString"`
}

func XmlInit() {
	testdata := Chat{Type: "Chat", ID: "asdf", Chat: "hansando"}

	file, err := os.Create("data.xml")
	if err != nil {
		fmt.Println("파일 생성 오류" , err)
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
	fmt.Println("XML Fin data.xml");
}
