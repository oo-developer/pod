package common

type Commands interface {
	Get(name string) (CommandService, bool)
}

type CommandService interface {
	Execute([]string) error
}
