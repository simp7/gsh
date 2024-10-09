package sh

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Env map[string]string

func (e Env) Set(token string) {
	key, value, ok := parseEnvToken(token)
	if ok {
		e[key] = value
	}
}

func (e Env) toArray() []string {
	result := make([]string, len(e))
	i := 0
	for key, value := range e {
		result[i] = key + "=" + value
		i++
	}
	return result
}

func (e Env) Unset(key string) {
	delete(e, key)
}

func (e Env) Get(key string) string {
	return e[key]
}

type sh struct {
	wd  string
	env Env
}

func New() (*sh, error) {
	s := new(sh)
	wd, err := os.Getwd()
	s.wd = wd
	s.env = make(Env)
	for _, token := range os.Environ() {
		s.env.Set(token)
	}
	return s, err
}

// export 사양
// 모든 큰따옴표 생략됨
// 띄어쓰기로 새 값 할당
// = 없는 경우 no-op(기존 값 그대로)
func parseEnvToken(token string) (string, string, bool) {
	parsed := strings.SplitN(token, "=", 2)
	if len(parsed) == 1 {
		return token, "", false
	}
	value := strings.ReplaceAll(parsed[1], `"`, "")
	return parsed[0], value, true
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
	newWd := target
	if target[0] != '/' {
		newWd = path.Join(s.wd, target)
	}

	if isDir(newWd) {
		s.wd = newWd
		return nil
	}
	return fmt.Errorf("invalid directory %s", newWd)
}

func (s *sh) pwd() string {
	return s.wd
}

func (s *sh) export(tokens []string) {
	for _, token := range tokens {
		s.env.Set(token)
	}
}

func (s *sh) envList() string {
	return strings.Join(s.env.toArray(), "\n")
}

func (s *sh) unset(keys []string) {
	for _, key := range keys {
		s.env.Unset(key)
	}
	return
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
		if len(seperated) == 1 {
			err = s.cd(s.env.Get("HOME"))
			break
		}
		err = s.cd(seperated[1])
	case "pwd":
		fmt.Println(s.pwd())
	case "export":
		if len(seperated) == 1 {
			fmt.Println(s.envList())
			break
		}
		s.export(seperated[1:])
	case "unset":
		if len(seperated) == 1 {
			err = errors.New("unset: not enough arguments")
			break
		}
		s.unset(seperated[1:])
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
	command.Env = s.env.toArray()

	output, err := command.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}
