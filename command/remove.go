package command

import "github.com/oo-developer/pod/common"

type removeCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (r *removeCommand) Execute(args []string) error {
	changeToPodWd(r.container, args...)
	r.container.Remove()
	return nil
}
