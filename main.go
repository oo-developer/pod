package main

import (
	"fmt"
	"os"

	"github.com/oo-developer/pod/command"
	"github.com/oo-developer/pod/config"
	"github.com/oo-developer/pod/container"
	"github.com/oo-developer/pod/system"
)

func main() {
	sys := system.Init()
	if sys.User() == "root" {
		fmt.Println("[ERROR] Pod can not run as root!")
		os.Exit(1)
	}
	conf := config.Init(sys)
	cont := container.Init(conf, sys)
	commands := command.Init(sys, conf, cont)
	if com, ok := commands.Get(os.Args[1]); ok {
		err := com.Execute(os.Args[2:])
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
		}
	} else {
		fmt.Println("[ERROR] unknown command!")
	}
}
