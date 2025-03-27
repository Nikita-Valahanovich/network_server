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

// Инициализация генератора случайных чисел
func init() {
	rand.Seed(time.Now().Unix())
}

// Пословицы
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

// Константы для настройки сервера
const (
	addr    = ":12345" // Порт, который будет слушать сервер
	network = "tcp4"   // Используемый сетевой протокол
)

func main() {
	// Создание TCP-сервера
	listen, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err)
	}
	// Гарантируем закрытие listener при завершении программы
	defer listen.Close()

	// Вывод сообщения о запуске сетевой службы
	fmt.Println("Go Proverbs Server is running on :12345")

	// Запускаем обработчик подключения в бесконечном цикле
	for {
		// Принимаем подключение
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue // Продолжаем работать при ошибках подключения
		}
		// Запускаем обработчик подключения в отдельной горутине
		// Это позволяет обрабатывать множество клиентов одновременно
		go handleConn(conn)
	}

}

// Обработчик подключения. Вызывается при каждом новом подключении.
func handleConn(conn net.Conn) {
	// закрытие соединения
	defer conn.Close()

	// Получаем адрес клиента для логирования
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("New connection from %s\r\n", clientAddr)

	// Канал для остановки отправки сообщений
	stopChan := make(chan struct{})

	// Горутина для чтения ввода от клиента
	go func() {
		reader := bufio.NewReader(conn)
		for {
			// Читаем данные до символа новой строки
			message, err := reader.ReadString('\n')
			if err != nil {
				// При ошибке чтения сигнализируем об остановке
				close(stopChan)
				return
			}
			// Проверяем, не пустой ли ввод (нажатие Enter)
			if strings.TrimSpace(message) == "" {
				fmt.Printf("Client %s requested disconnect\n", clientAddr)
				close(stopChan)
				return
			}
		}
	}()
	// Создаем тикер для отправки сообщений каждые 3 секунды
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop() // Гарантируем освобождение ресурсов

	for {
		select {
		case <-ticker.C:
			// При срабатывании тикера выбираем случайную поговорку
			proverb := proverbs[rand.Intn(len(proverbs))]

			// Отправляем поговорку клиенту
			_, err := conn.Write([]byte(proverb + "\r\n"))
			if err != nil {
				// При ошибке записи закрываем соединение
				fmt.Printf("Client %s disconnected: %v\n", clientAddr, err)
				return
			}
		case <-stopChan:
			// Получен сигнал остановки (клиент нажал Enter или разорвал соединение)
			fmt.Printf("Closing connection with %s\n", clientAddr)
			return
		}
	}
}
