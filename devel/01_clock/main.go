package main

import (
	"fmt"
	"log"
)

func main() {
	c, err := wbtimer.New(wbtimer.DefaultHost)

	if err != nil {
		log.Fatalln("Error!", err)
	}

	fmt.Println(c.CurrentTime())
}
