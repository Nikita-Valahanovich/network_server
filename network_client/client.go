package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// подключение к сетевой службе
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	// закрываем ресурс
	defer conn.Close()

	// Добавляем входящее сообщение в буфер
	reader := bufio.NewReader(conn)
	// Читаем сообщение из буфера
	fmt.Println(reader)
}
