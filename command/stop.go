package command

import "github.com/oo-developer/pod/common"

type stopCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *stopCommand) Execute(strings []string) error {
	s.container.Stop()
	return nil
}
