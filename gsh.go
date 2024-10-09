package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/simp7/gsh/sh"
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
	s, err := sh.New()
	if err != nil {
		fmt.Printf("gsh start failed: %s\n", err)
		os.Exit(1)
		return
	}

	for {
		fmt.Print("gsh> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			exitWithError(err)
			return
		}
		input = strings.TrimSpace(input)

		if err = s.Execute(input); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Println("No such file or directory (os error 2)")
				continue
			}
			fmt.Println(err)
		}
	}
}
