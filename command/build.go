package command

import "github.com/oo-developer/pod/common"

type buildCommand struct {
	sys       common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (b *buildCommand) Execute(strings []string) error {
	b.container.BuildDockerFile()
	return nil
}
