package command

import "github.com/oo-developer/pod/common"

type listCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (l listCommand) Execute(strings []string) error {
	l.container.ListPods()
	return nil
}
