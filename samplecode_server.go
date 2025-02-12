package main

import (
	"fmt"
	"net"
	"sync"
	"io"
)

// Client 구조체는 각 연결된 클라이언트의 정보를 저장합니다
type Client struct {
	conn     net.Conn
	id       string
	outbound chan []byte
}

// Server 구조체는 서버의 상태를 관리합니다
type Server struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}


// NewServer는 새로운 서버 인스턴스를 생성합니다
func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start는 서버의 메인 루프를 시작합니다
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
				default:
					close(client.outbound)
					delete(s.clients, client.id)
				}
			}
			s.mu.RUnlock()
		}
	}

}

// handleClient는 각 클라이언트의 연결을 처리합니다
func (s *Server) handleClient(client *Client) {
	defer func() {
		s.unregister <- client
		client.conn.Close()
	}()

	// 클라이언트로부터 메시지를 읽는 고루틴
	go func() {
		for {
			recv := make([]byte, 4096)
			n, err := client.conn.Read(recv)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("Client disconnected: %s\n", client.id)
					return
				}
				fmt.Printf("Error reading from client: %v\n", err)
				return
			}
			if n > 0 {
				// 받은 메시지를 모든 클라이언트에게 브로드캐스트
				s.broadcast <- recv[:n]
			}
		}
	}()

	// 클라이언트에게 메시지를 보내는 루프
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
	// 서버 메인 루프 시작
	go server.Start()

	// TCP 리스너 설정
	l, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		fmt.Println("Failed to Listen:", err)
		return
	}
	defer l.Close()

	fmt.Println("Server started on port 8000")

	// 클라이언트 연결 수락 루프
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept: %v\n", err)
			continue
		}

		client := &Client{
			conn:     conn,
			id:       conn.RemoteAddr().String(),
			outbound: make(chan []byte, 100),
		}

		server.register <- client
		go server.handleClient(client)
	}
}
