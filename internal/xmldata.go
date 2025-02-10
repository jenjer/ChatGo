package xmldata

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Data struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
	City    string   `xml:"city"`
}

func XmlInit() {
	testdata := Data{Name: "leesonsin", Age: 45, City: "hansando"}

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
