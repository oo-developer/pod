package command

import "github.com/oo-developer/pod/common"

type runCommand struct {
	sys       common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (r runCommand) Execute(args []string) error {
	changeToPodWd(r.container, args...)
	r.container.RunContainer()
	return nil
}
