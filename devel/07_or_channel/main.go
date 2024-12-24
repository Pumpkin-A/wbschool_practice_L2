package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("fone after %v\n", time.Since(start))
}

func or(chans ...<-chan interface{}) <-chan interface{} {
	firstValueRecieved := make(chan interface{})

	waitUntil := func(waitForClose <-chan interface{}) {
		for {
			select {
			case _, more := <-firstValueRecieved:
				if !more {
					return
				}

			case _, more := <-waitForClose:
				if !more {
					close(firstValueRecieved)
					return
				}
			}
		}
	}

	for i := 0; i < len(chans); i++ {
		go waitUntil(chans[i])
	}

	return firstValueRecieved
}
