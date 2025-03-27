package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var proverbs = []string{
	"Don't communicate by sharing memory, share memory by communicating.",
	"Concurrency is not parallelism.",
	"Channels orchestrate; mutexes serialize.",
	"The bigger the interface, the weaker the abstraction.",
	"Make the zero value useful.",
	"interface{} says nothing.",
	"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
	"A little copying is better than a little dependency.",
	"Syscall must always be guarded with build tags.",
	"Cgo must always be guarded with build tags.",
	"Cgo is not Go.",
	"With the unsafe package there are no guarantees.",
	"Clear is better than clever.",
	"Reflection is never clear.",
	"Errors are values.",
	"Don't just check errors, handle them gracefully.",
	"Design the architecture, name the components, document the details.",
	"Documentation is for users.",
	"Don't panic.",
}

const addr = ":12345"
const network = "tcp4"

func main() {
	listen, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()
	fmt.Println("Go Proverbs Server is running on :12345")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("New connection from %s\r\n", clientAddr)

	stopChan := make(chan struct{})

	// Горутина для чтения ввода от клиента
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				close(stopChan)
				return
			}
			// Если получен пустой ввод (просто Enter)
			if strings.TrimSpace(message) == "" {
				fmt.Printf("Client %s requested disconnect\n", clientAddr)
				close(stopChan)
				return
			}
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proverb := proverbs[rand.Intn(len(proverbs))]
			_, err := conn.Write([]byte(proverb + "\r\n"))
			if err != nil {
				fmt.Printf("Client %s disconnected: %v\n", clientAddr, err)
				return
			}
		case <-stopChan:
			fmt.Printf("Closing connection with %s\n", clientAddr)
			return
		}
	}
}
