package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getDataFromUser() (string, time.Duration) {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout in seconds")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatalf("usage: go-telnet --timeout=10s host port")
	}
	host := args[0]
	port := args[1]

	return host + ":" + port, *timeout
}

func main() {
	// Инициализируем нашего клиента переденными флагами из командной строки
	client := NewTelnetClient(getDataFromUser())
	if err := client.initConnection(); err != nil {
		log.Fatalf("error while connecting %v", err)
	}
	// Откладываем закрытие соединения
	defer func() {
		if err := client.closeConnection(); err != nil {
			log.Fatalf("error while closing conn")
		}
		fmt.Println("conn socket closed")
	}()

	// Создаем перехватчк завершения работы программы
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	errors := make(chan error, 1)

	// Создаем горутину для чтения данных из сокета соединения telnet
	go func() {
		for {
			err := client.receieveMsg()
			if err != nil {
				errors <- err
				return
			}
		}
		// В случае ошибок - отправляем ошибку в канал ошибок и выходим из горутины
	}()

	// Создаем горутину для чтения данных из консоли (STDIN)
	go func() {
		// Обрабатываем ввод данных из STDIN до нажатия enter
		inputReader := bufio.NewReader(os.Stdin)
		for {
			line, err := inputReader.ReadString('\n')
			if err != nil {
				errors <- err
				return
			}
			if err := client.sendMsg(line); err != nil {
				errors <- err
				return
			}
		}
		// В случае ошибок - отправляем ошибку в канал ошибок и выходим из горутины
	}()

	// Запускаем select в цикле, который отслеживает завершающие сигналы от ОС и наличие ошибок в канале
	// При получения любого из них завершает программу, закрывая сокет соединения с сервером в верхнем defer
	for {
		select {
		case <-signals:
			return
		case err := <-errors:
			fmt.Println("Error from client: ", err.Error())
			if err != nil {
				return
			}
		default:
			continue
		}
	}
}
