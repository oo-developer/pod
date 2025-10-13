package command

import "github.com/oo-developer/pod/common"

type startCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *startCommand) Execute(args []string) error {
	changeToPodWd(s.container, args...)
	s.container.Start()
	return nil
}
