package command

import "github.com/oo-developer/pod/common"

type recipesCommand struct {
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func (l *recipesCommand) Execute(strings []string) error {
	if len(strings) == 0 {
		l.config.ListRecipes()
		return nil
	}
	switch strings[0] {
	case "list":
		l.config.ListRecipes()
		return nil
	case "update":
		break
	}
	return nil
}
