package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleIncomingMessage(conn net.Conn) {
	defer wg.Done()
	for {
		data := make([]byte, 1024)
		conn.Read(data)
		fmt.Println(string(data))
		if string(data) == "" {
			break
		}
		time.Sleep(time.Second * 1)
	}

}

func handleUserInput(conn net.Conn) {
	inputStream := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message:")
		inputStream.Scan()
		text := inputStream.Text()
		if text == "Q" {
			conn.Write([]byte(text))
			break
		}
		conn.Write([]byte(text))

	}
	conn.Close()
	wg.Done()
	wg.Done()
}
func main() {

	fmt.Println("Starting Client")
	connection, err := net.Dial("tcp", "localhost:8000")
	handleError(err)
	defer connection.Close()
	//connection.Write([]byte("Aravind"))
	wg.Add(1)
	go handleIncomingMessage(connection)
	wg.Add(1)
	go handleUserInput(connection)
	wg.Wait()

}
