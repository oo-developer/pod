package command

import "github.com/oo-developer/pod/common"

type statusCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *statusCommand) Execute(strings []string) error {
	s.container.Status()
	return nil
}
