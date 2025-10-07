package command

import "github.com/oo-developer/pod/common"

type shellCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (s *shellCommand) Execute(args []string) error {

	if common.HasOption("--script", args...) {
		s.container.BuildShellScript()
		return nil
	}
	s.container.Shell()
	return nil
}
