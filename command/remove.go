package command

import "github.com/oo-developer/pod/common"

type removeCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (r *removeCommand) Execute(strings []string) error {
	r.container.Remove()
	return nil
}
