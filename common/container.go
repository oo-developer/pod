package common

type ContainerService interface {
	WriteDefaultPod()
	GetDefaultPod() *PodDefinition
	BuildDockerFile()
	RunContainer()
	Start()
	Stop()
	Remove()
	Shell()
	BuildShellScript()
	Status()
}
