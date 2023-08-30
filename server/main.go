package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var connectionsMutex sync.Mutex

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var connections []net.Conn

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer wg.Done()
	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		handleError(err)
		if string(data)[:1] == "Q" {
			fmt.Println("Here")
			fmt.Printf("Connection %s is leaving....Good Bye\n", conn.RemoteAddr())
			connectionsMutex.Lock()
			// Remove the disconnected connection from the slice
			for i, c := range connections {
				if c == conn {
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			connectionsMutex.Unlock()

			break
		}
		fmt.Printf("Message [CLIENT %s]:  %s\n", conn.RemoteAddr(), string(data))
		connectionsMutex.Lock()
		for i := 0; i < len(connections); i++ {
			if connections[i] != conn {
				connections[i].Write([]byte(fmt.Sprintf("\nMESSAGE FROM [CLIENT %s]: %s\n", conn.RemoteAddr(), strings.TrimSpace(string(data)))))
			}
		}
		connectionsMutex.Unlock()
	}
}
func main() {

	fmt.Println("Starting Server......")
	listener, err := net.Listen("tcp", ":8000")
	handleError(err)
	defer listener.Close()
	defer wg.Wait()
	for {
		conn, err := listener.Accept()
		handleError(err)
		fmt.Printf("[CLIENT %s ] Connected...\n", conn.RemoteAddr())
		connectionsMutex.Lock()
		connections = append(connections, conn)
		connectionsMutex.Unlock()
		wg.Add(1)
		go handleConnection(conn)
	}
}
