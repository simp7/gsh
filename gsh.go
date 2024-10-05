package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func execute(input string) error {
	seperated := strings.Split(input, " ")
	args := seperated[1:]

	command := exec.Command(seperated[0], args...)
	output, err := command.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("gsh> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			exitWithError(err)
			return
		}
		input = strings.TrimSpace(input)

		if input == "exit" {
			return
		}

		if err := execute(input); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Println("No such file or directory (os error 2)")
				continue
			}
			exitWithError(err)
			return
		}
	}
}
