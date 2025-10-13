package command

import (
	"strings"

	"github.com/oo-developer/pod/common"
)

type commands struct {
	commands  map[string]common.CommandService
	system    common.SystemService
	config    common.ConfigService
	container common.ContainerService
}

func Init(system common.SystemService, config common.ConfigService, container common.ContainerService) common.Commands {
	c := &commands{
		commands:  make(map[string]common.CommandService),
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["shell"] = &shellCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["init"] = &initCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["build"] = &buildCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["run"] = &runCommand{
		sys:       system,
		config:    config,
		container: container,
	}
	c.commands["stop"] = &stopCommand{
		container: container,
		config:    config,
		system:    system,
	}
	c.commands["start"] = &startCommand{
		container: container,
		config:    config,
		system:    system,
	}
	c.commands["recipes"] = &recipesCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["remove"] = &removeCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["status"] = &statusCommand{
		system:    system,
		config:    config,
		container: container,
	}
	c.commands["list"] = &listCommand{
		system:    system,
		config:    config,
		container: container,
	}
	return c
}

func (c *commands) Get(name string) (common.CommandService, bool) {
	command, ok := c.commands[name]
	return command, ok
}

func hasPodName(args ...string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return true
		}
	}
	return false
}

func getPodName(args ...string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}
	return ""
}

func changeToPodWd(container common.ContainerService, args ...string) {
	if hasPodName(args...) {
		name := getPodName(args...)
		container.ChangeToPodDir(name)
	}
}
