package sh

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type sh struct {
	wd string
}

func New() (*sh, error) {
	s := new(sh)
	wd, err := os.Getwd()
	s.wd = wd
	return s, err
}

func (s *sh) exit() {
	os.Exit(0)
}

func isDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (s *sh) cd(target string) error {
	newWd := path.Join(s.wd, target)
	if isDir(newWd) {
		s.wd = newWd
		return nil
	}
	return fmt.Errorf("invalid directory %s", newWd)
}

func (s *sh) pwd() string {
	return s.wd
}

func (s *sh) executeBuiltIn(seperated []string) (bool, error) {
	var err error
	if len(seperated) == 0 || seperated[0] == "" {
		return true, nil
	}

	switch seperated[0] {
	case "exit":
		s.exit()
	case "cd":
		if len(seperated) != 2 {
			return true, errors.New("command 'cd' should contains one argument")
		}
		err = s.cd(seperated[1])
	case "pwd":
		fmt.Println(s.pwd())
	default:
		return false, err
	}
	return true, err
}

func (s *sh) Execute(input string) error {
	seperated := strings.Split(input, " ")
	args := seperated[1:]

	if isBuiltin, err := s.executeBuiltIn(seperated); isBuiltin {
		return err
	}

	command := exec.Command(seperated[0], args...)
	command.Dir = s.wd

	output, err := command.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}
