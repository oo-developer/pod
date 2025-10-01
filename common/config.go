package common

type ConfigService interface {
	BasePath() string
	LibraryPath() string
	PodsPath() string
	PodPath(folder string) string
	PodImageName(folder string) string
	ConfigPath() string
	LibraryGitRepository() string
	AuthorizedKey() string
	PrivateKeyPath() string
}
