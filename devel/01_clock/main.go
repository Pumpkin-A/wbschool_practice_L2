package main

import (
	"fmt"
	"log"

	"github.com/Pumpkin-A/wbschool_practice_L2/devel/01_clock/wbtimer"
)

func main() {
	c, err := wbtimer.New(wbtimer.DefaultHost)

	if err != nil {
		log.Fatalln("Error!", err)
	}

	fmt.Println(c.CurrentTime())
}
