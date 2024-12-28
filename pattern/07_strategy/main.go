package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Паттерн "стратегия" реализован на примере динамического выбора конкретного хранилища кэша.
// Стратегия позволяет варьировать поведение объекта во время выполнения программы, подставляя в него различные
// объекты-поведения. Интерфейс CachingStrategy определяет общее поведение, которое реализует каждое конкретное хранилище,
// благодаря этому не приходится хранить все возможные реализации в одном классе, а выбирать только нужную в процессе работы программы

// Плюсы:  Горячая замена алгоритмов на лету.
// Изолирует код и данные алгоритмов от остальных классов.
// Уход от наследования к делегированию.
// Реализует принцип открытости/закрытости.
// Минусы: Усложняет программу за счёт дополнительных классов.
// Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.

type CachingStrategy interface {
	connect() error
	write(data string)
}

type Redis struct {
}

func (r *Redis) connect() error {
	fmt.Println("выбрана стратегия кэширования в redis")
	return nil
}

func (r *Redis) write(data string) {
	fmt.Printf("запись кэша в redis: %s\n", data)
}

type InMemory struct {
}

func (m *InMemory) connect() error {
	fmt.Println("выбрана стратегия кэширования в мапу")
	return nil
}

func (r *InMemory) write(data string) {
	fmt.Printf("запись кэша в мапу: %s\n", data)
}

type File struct {
}

func (r *File) connect() error {
	// fmt.Println("выбрана стратегия кэширования в текстовый документ")
	return errors.New("any error")
}

func (r *File) write(data string) {
	fmt.Printf("запись кэша в файл: %s\n", data)
}

type CacheStorage struct {
	// конкретный storage здесь не задается, это должно быть реализовано внутри каждой стратегии
	cachingStrategy CachingStrategy
	capacity        int
	maxCapacity     int
}

func initCacheStorage(s CachingStrategy) *CacheStorage {
	err := s.connect()
	if err != nil {
		fmt.Println("проблема подключения новой стратегии кэша, стратегия не установлена")
		return nil
	}
	return &CacheStorage{
		cachingStrategy: s,
		capacity:        0,
		maxCapacity:     2,
	}
}

func (c *CacheStorage) setCachingStrategy(s CachingStrategy) {
	err := s.connect()
	if err != nil {
		fmt.Println("проблема подключения новой стратегии кэша, стратегия не изменена")
		return
	}
	c.cachingStrategy = s
}

func main() {
	inMemory := &InMemory{}
	cacheStorage := initCacheStorage(inMemory)
	redis := &Redis{}
	file := &File{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			cacheStorage.cachingStrategy.write("hello cache")
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			i := rand.Intn(10)
			if i%3 == 0 {
				cacheStorage.setCachingStrategy(redis)
			} else if i%3 == 1 {
				cacheStorage.setCachingStrategy(file)
			} else {
				cacheStorage.setCachingStrategy(inMemory)
			}
			time.Sleep(time.Second * 3)
		}
	}()
	wg.Wait()
}
