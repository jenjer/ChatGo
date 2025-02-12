package main

import (
	"fmt"
	"io"
	"net"
	//"encoding/xml"
	//xmlstruct "github.com/jenjer/ChatGo/internal"
	"sync"
	loginModule "github.com/jenjer/ChatGo/internal/serverPackage/Login"
)

type Client struct {
	conn		net.Conn
	id			string
	outbound	chan []byte
}

type Server struct {
	clients		map[string] *Client
	mu			sync.RWMutex
	broadcast	chan []byte
	register	chan *Client
	unregister	chan *Client
}

func NewServer() *Server {
	return &Server {
		clients:	make(map[string] *Client),
		broadcast:	make(chan []byte),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
	}
}

func (s *Server) Start(){
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept: %v\n", err)
			continue
		}

		// 로그인 시도를 성공하면 하단으로 가고 아니면 그냥 패스해야됨
		//개별적으로 close 해야될것들은 handleClient 에서 처리
		if (loginModule.TryLogin(conn) == true) {
			client := &Client {
				conn:		conn,
				id:			conn.RemoteAddr().String(),
				outbound:	make(chan[]byte, 100),
			}
			server.register <- client
			go server.handleClient(client)
		}
	}
}

