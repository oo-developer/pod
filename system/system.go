package system

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/oo-developer/pod/common"
	"golang.org/x/term"
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

func (s *system) ExecuteShell(command string, args ...string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting raw mode: %v\n", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {

			fmt.Printf("[ERROR] %s\n", exitError.Stderr)
		}
		os.Exit(1)
	}
	go func() {
		_, err := io.Copy(os.Stdout, stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdout: %v\n", err)
		}
	}()
	go func() {
		_, err := io.Copy(os.Stderr, stderr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stderr: %v\n", err)
		}
	}()
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				}
				stdin.Close()
				return
			}
			if n == 0 {
				continue
			}
			if buf[0] == 4 { // Ctrl+D
				stdin.Close()
				return
			}
			_, err = stdin.Write(buf[:n])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to stdin: %v\n", err)
				stdin.Close()
				return
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Command finished with error: %v\n", err)
	}
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
