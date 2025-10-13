package command

import "github.com/oo-developer/pod/common"

type statusCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *statusCommand) Execute(args []string) error {
	changeToPodWd(s.container, args...)
	s.container.Status()
	return nil
}
