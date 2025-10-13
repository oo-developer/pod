package command

import (
	"fmt"
	"os"

	"github.com/oo-developer/pod/common"
)

type recipesCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (l *recipesCommand) Execute(args []string) error {
	if len(args) < 1 {
		fmt.Println("[ERROR] No os flavor specified. Please use: pod recipes <os flavor> (e.g. debian)")
		os.Exit(1)
	}
	l.config.ListRecipes(args[0])
	return nil
}
