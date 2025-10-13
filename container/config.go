package container

type Config struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Path       string `json:"path"`
	BaseImage  string `json:"baseImage"`
	Flavor     string `json:"flavor"`
	FirstStart bool   `json:"firstStart"`
}
