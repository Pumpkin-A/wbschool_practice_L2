package main

import (
	"fmt"
	"time"
)

func main() {
	// Создаём функцию-констркутор каналов, которые закроются через определённое время
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			c <- 1 // имитируем работу каналов: они не просто закрываются, но и могут прислать данные
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

func or(chans ...<-chan interface{}) <-chan interface{} {
	// Создаем новый канал, в который будем сливать любое количество пришедших каналов
	orChan := make(chan interface{})

	// Запускаем горутины дял каждого канала, которые ожидают их закрытия
	for i := 0; i < len(chans); i++ {
		go func(waitForClose <-chan interface{}) {
			for {
				select {
				// При закрытии нового канала orChan - завершаем горутину, чтобы освободить ресурсы
				case _, isOpen := <-orChan:
					if !isOpen {
						return
					}
				// При закрытия любого из них - закрываем новый канал
				case _, isOpen := <-waitForClose:
					if !isOpen {
						close(orChan)
						return
					}
				}
			}
		}(chans[i])
	}

	return orChan
}
