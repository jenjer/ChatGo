package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"sync"

	xmlstruct "github.com/jenjer/ChatGo/internal"
	define "github.com/jenjer/ChatGo/internal/clientPackage/defines"
	ini "github.com/jenjer/ChatGo/internal/clientPackage/iniFunc"
	login "github.com/jenjer/ChatGo/internal/clientPackage/login"
	//Global "github.com/jenjer/ChatGo/internal/clientPackage"
)

func makeXml(input string) xmlstruct.Chat {
	var sendChat xmlstruct.Chat
	sendChat.Type = "Chat"
	sendChat.ID = "input ID"
	sendChat.Chat = input
	return sendChat
}

func CheckLogin(bytes []byte) bool{
	var Ans  xmlstruct.LoginAns
	err := xml.Unmarshal(bytes, &Ans)
	if err != nil {
		fmt.Println("Error parsing LogAnsXml: ", err)
		return false
	}
	fmt.Println("login result : ", Ans.Result)
	return true
}

func makeString(bytes []byte) string {
	var msg xmlstruct.Chat
	err := xml.Unmarshal(bytes, &msg)
	if err != nil {
		if (CheckLogin(bytes) == true){
			return ""
		} else {
			fmt.Println("Error parsing ChatXml: ", err)
			return ""
		}
	}
	return msg.ID + " : " + msg.Chat
}

func iniMain() (string, bool) {
	args := os.Args
	ip := ini.GetIni(define.MainData, define.ServerIP)
	if len(os.Args) == 1 {

		/*
			propertySection := cfg.Section("property")
			ServerIP := propertySection.Key("ServerIP").String()
		*/
		//여기 까지 ini 파일에서 server ip 주소 받아오는데 성공했는데 없다면?
		if ip == "" {
			fmt.Println("default from ini file is null")
			fmt.Println("Please reload program with ServerIP")
			return "", false
		} else {
			fmt.Println("server ip : " + ip)
			return ip, true
		}
	} else {
		ip = args[1]
		fmt.Printf("is this server ip? (%s)\n(y/n) :", ip)
		var temp string
		fmt.Scanln(&temp)
		if temp != "y" {
			return "", false
		}
		return ip, true
	}
}

func main() {
	ip, errtemp := iniMain()
	if errtemp == false {
		return
	}

	xmlstruct.XmlInit()
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("failed to dial:", err)
		return
	}
	defer conn.Close()

	//login with Conn
	login.Login(conn)
	var wg sync.WaitGroup
	wg.Add(2)

	scanner := bufio.NewScanner(os.Stdin)
	// 데이터 보내는 goroutine
	go func(c net.Conn) {
		defer wg.Done()
		defer c.Close()
		fmt.Print("> ")
		for {
			var input string
			scanner.Scan()
			input = scanner.Text()
			//fmt.Scanln(&input)

			if input == "quit" {
				fmt.Println("Closing connection...")
				return
			}

			if input == "" {
				fmt.Print("\r") // string(recv[:n]))
				fmt.Print(">")
				continue
			}

			Chatstruct := makeXml(input)
			encoder := xml.NewEncoder(c)
			err := encoder.Encode(Chatstruct)
			//_, err = c.Write([]byte(input))
			if err != nil {
				fmt.Println("Failed to write data:", err)
				return
			}
		}
	}(conn)

	// 데이터 받는 goroutine
	go func(c net.Conn) {
		defer wg.Done()
		defer c.Close()
		recv := make([]byte, 4096)
		for {
			_, err := c.Read(recv)
			if err != nil {
				fmt.Println("Failed to read data:", err)
				return
			}
			forPrint := makeString(recv)
			fmt.Print("\r" + forPrint + "\n") // string(recv[:n]))
			fmt.Print("> ")
		}
	}(conn)

	wg.Wait()
}
