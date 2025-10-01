package system

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/oo-developer/pod/common"
)

type system struct {
}

func Init() common.SystemService {
	return &system{}
}

func (s *system) Execute(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("[ERROR] %s\n", exitError.Stderr)
		}
		os.Exit(1)
	}
	return string(output)
}

func (s *system) User() string {
	return os.Getenv("USER")
}

func (s *system) Uid() int {
	return os.Getuid()
}

func (s *system) Gid() int {
	return os.Getgid()
}

func (s *system) HomeDir() string {
	return os.Getenv("HOME")
}

func (s *system) Display() string {
	return os.Getenv("DISPLAY")
}

func (s *system) Architecture() string {
	return strings.Trim(s.Execute("uname", "-m"), "\n")
}

func (s *system) CheckOrCreatePath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0750)
			if err != nil {
				fmt.Printf("[ERROR] %v\n", err)
				os.Exit(1)
			}
		}
		return false
	}
	return true
}

func (s *system) PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (s *system) FreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return port
}

func (s *system) CopyFile(src, target string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}
	err = os.WriteFile(target, data, srcInfo.Mode())
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return err
	}
	return nil
}
