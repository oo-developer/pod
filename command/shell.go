package command

import "github.com/oo-developer/pod/common"

type shellCommand struct {
	sys       common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *shellCommand) Execute(strings []string) error {
	return nil
}
