package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func printUsageHelp() {
	log.Printf("Commands: cd, pwd, echo, kill, ps, exit")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	printUsageHelp()
	os.Exit(exitcode)
}

func main() {
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = printUsageHelp
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	// Печатем приглашение к вводу и ожидаем строчки от пользователя
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\r')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// На винде небходимо чистить ввод от символов
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")
	input = strings.TrimPrefix(input, "\r")
	input = strings.TrimPrefix(input, "\n")

	args := strings.Split(input, " ")
	fmt.Println(args)

	// Команду изменения директории необходимо спарсить и преобразовать
	switch args[0] {
	case "cd":
		if len(args) != 2 {
			return errors.New("one directory path is required")
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
