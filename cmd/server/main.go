package main

import (
	"fmt"
	"io"
	"net"
	"net/http"

	//"encoding/xml"
	"encoding/xml"
	"sync"

	xmlstruct "github.com/jenjer/ChatGo/internal"
	DBConn "github.com/jenjer/ChatGo/internal/serverPackage/DB"
	loginModule "github.com/jenjer/ChatGo/internal/serverPackage/Login"
	"github.com/labstack/echo/v4"
)

type Client struct {
	conn     net.Conn
	id       string
	outbound chan []byte
	loginid  string
}

type Server struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func messageStruct(message []byte) string {
	var msg xmlstruct.Chat
	err := xml.Unmarshal(message, &msg)
	if err != nil {
		fmt.Println("Error parsing XmL:", err)
		return ""
	}
	fmt.Printf("\nParsingData\n")
	fmt.Printf("Message Type : %s\n", msg.Type)
	fmt.Printf("Message ID : %s\n", msg.ID)
	fmt.Printf("Message Chatstring : %s\n", msg.Chat)
	return msg.ID + " : " + msg.Chat
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.id] = client
			s.mu.Unlock()
			fmt.Printf("New client registered: %s\n", client.id)
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.outbound)
			}
			s.mu.Unlock()
			fmt.Printf("Client unregistered: %s\n", client.id)
		case message := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.outbound <- message:
				default: //기본 100개의 채팅을 받을 수 있는데 100 개를 다 못치는 상태이기 때문에 그냥 종료해버린다.
					close(client.outbound)
					delete(s.clients, client.id)
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *Server) handleClient(client *Client) {
	defer func() {
		s.unregister <- client
		client.conn.Close()
	}()

	go func() {
		for {
			recv := make([]byte, 4096)
			n, err := client.conn.Read(recv)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("Client disconnected: %s\n", client.id)
					return
				}
				fmt.Printf("Error reading from clinet: %v\n", err)
				return
			}
			if n > 0 {
				s.broadcast <- recv[:n]
			}
		}
	}()

	for message := range client.outbound {
		_, err := client.conn.Write(message)
		if err != nil {
			fmt.Printf("Error writing to client: %v\n", err)
			return
		}
	}
}
func loginSuccessMessage(client *Client, conn net.Conn, success bool) {
	message := []byte("success")
	_, err := client.conn.Write(message)
	if err != nil {
		fmt.Printf("Error writing to client: %v\n", err)
		return
	}
	///////
	var LoginAnswer string
	if success == true {
		LoginAnswer = "Success"
	} else {
		LoginAnswer = "False"
	}
	var senddata xmlstruct.LoginAns
	senddata.Type = "LoginResult"
	senddata.Result = LoginAnswer
	encoder := xml.NewEncoder(conn)
	err = encoder.Encode(senddata)
	if err != nil {
		fmt.Println("Failed to send data loginfail")
		return
	}
}

func main() {
	server := NewServer()

	go server.Start()

	l, err := net.Listen("tcp", "0.0.0.0:8000")

	if err != nil {
		fmt.Println("Failed to Listen: ", err)
		return
	}
	defer l.Close()

	fmt.Println("Server started on port 8000")
	// 현 위치에 고 루틴으로 웹서버를 넣을것
	go func() {
		e := echo.New()
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
		e.Logger.Fatal(e.Start(":1323"))
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept: %v\n", err)
			continue
		}

		// 로그인 시도를 성공하면 하단으로 가고 아니면 그냥 패스해야됨
		//개별적으로 close 해야될것들은 handleClient 에서 처리
		DbConn, terr := DBConn.NewUserDB("./testdb")
		if terr != nil {
			fmt.Printf("err : %v\n", terr)
		}
		testerr := DbConn.AddUser("asdf", "asdf")
		testerr = DbConn.AddUser("first", "first")
		testerr = DbConn.AddUser("test", "1234")
		client := &Client{
			conn:     conn,
			id:       conn.RemoteAddr().String(),
			loginid:  "",
			outbound: make(chan []byte, 100),
		}
		if testerr != nil {
			fmt.Printf("err : %v\n", testerr)
			if booltemp, recvID := loginModule.TryLogin(conn, DbConn); booltemp == true {
				server.register <- client
				client.loginid = recvID
				//send login success message
				loginSuccessMessage(client, conn, true)
				go server.handleClient(client)
				fmt.Println("login Success")
			} else {
				fmt.Println("login failed")
				loginSuccessMessage(client, conn, false)
				// login fail message send
			}
		}
	}
}
