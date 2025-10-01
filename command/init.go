package command

import (
	"github.com/oo-developer/pod/common"
)

type initCommand struct {
	sys       common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (i *initCommand) Execute(strings []string) error {
	i.container.WriteDefaultPod()
	return nil
}
