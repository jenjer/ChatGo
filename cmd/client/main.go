package main

import(
	myxml "github.com/jenjer/ChatGo/internal"
	"os"
	ini "github.com/jenjer/ChatGo/internal/clientPackage/iniFunc"
	"fmt"
	define "github.com/jenjer/ChatGo/internal/clientPackage/defines"
	login "github.com/jenjer/ChatGo/internal/clientPackage/login"
	"net"
	"sync"
	"time"
)

func iniMain()(string, bool){
	args := os.Args
	ip := ini.GetIni(define.MainData, define.ServerIP);
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
		return ip,true
	}
}

func main() {

	time.Sleep(10000)
	ip, errtemp := iniMain()
	if errtemp == false {
		return
	}

	myxml.XmlInit()
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("failed to dial:", err)
		return
	}
	defer conn.Close()

	login.Login()

	login.set
	var wg sync.WaitGroup
	wg.Add(2)

	// 데이터 보내는 goroutine
	go func(c net.Conn) {
		defer wg.Done()
		defer c.Close()
		for {
			var input string
			fmt.Print("MyChat :  ")
			fmt.Scanln(&input)

			if input == "quit" {
				fmt.Println("Closing connection...")
				return
			}

			if input == "" {
				continue
			}

			_, err = c.Write([]byte(input))
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
			n, err := c.Read(recv)
			if err != nil {
				fmt.Println("Failed to read data:", err)
				return
			}
			fmt.Printf("\nServer: %s\n", string(recv[:n]))
		}
	}(conn)

	wg.Wait()
}
