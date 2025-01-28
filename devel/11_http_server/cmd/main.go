package main

import (
	"calendar-server/config"
	eventstorage "calendar-server/eventStorage"
	"calendar-server/filedb"
	"calendar-server/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.NewDefaultConfig()

	db, err := filedb.New(cfg.DbFilename)
	if err != nil {
		log.Fatalf("NewFileDB: %s", err.Error())
	}

	es, err := eventstorage.New(db)
	if err != nil {
		log.Fatalf("eventstorage: %s", err.Error())
	}

	server := server.New(*cfg, es)

	// Создаем каналы для перехвата сигнала от ОС и синхронизации закрытия сервера.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server: started at address %v", cfg.Port)
		serverErrors <- http.ListenAndServe(cfg.Port, nil)
	}()

	exitCh := make(chan struct{})
	// В данной горутину ждём ошибку от сервера или сигнал от ОС и далее стартуем завершение работы, отправляя сообщение в exitCh
	go func() {
		defer func() {
			exitCh <- struct{}{}
		}()

		for {
			select {
			case err := <-serverErrors:
				fmt.Printf("server error %v\n", err)
				return
			case <-osSignals:
				fmt.Println("service shutdown...")
				return
			default:
				continue
			}
		}
	}()

	<-exitCh

	// Запускаем деструкторы для наших сущностей, чтобы они корректно завершили работу.

	// Сначала ожидаем закрытия сервера, чтобы текущие запросы окончили работу. Новые не принимаются
	if err := server.Shutdown(); err != nil {
		log.Printf("Error while closing server; %s", err.Error())
	} else {
		log.Println("Server closed!")
	}

	// Закрываем слой хранения событий, чтобы изменения записались в файл-БД
	if err := es.Close(); err != nil {
		log.Printf("Error while ES close; %s", err.Error())
	} else {
		log.Println("ES closed!")
	}

	// Коррктено отдаём ресурс - закрываем файловый дескриптор
	if err := db.Close(); err != nil {
		log.Printf("Error while fileDb closing; %s", err.Error())
	} else {
		log.Println("DB closed!")
	}
}
