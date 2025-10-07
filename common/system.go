package common

type SystemService interface {
	Execute(command string, args ...string) string
	ExecuteShell(command string, args ...string)
	User() string
	Uid() int
	Gid() int
	HomeDir() string
	Display() string
	Architecture() string
	CheckOrCreatePath(path string) bool
	PathExists(path string) bool
	FreePort() int
	CopyFile(src, target string) error
}
