package common

type Container struct {
	Image  string `json:"image"`
	Flavor string `json:"flavor"`
}

type Ssh struct {
	PrivateKeyPath string `json:"privateKeyPath"`
	AuthorizedKey  string `json:"authorizedKey"`
}

type PodDefinition struct {
	Container Container `json:"container"`
	Ssh       Ssh       `json:"ssh"`
	Packages  []string  `json:"packages"`
	Recipes   []string  `json:"recipes"`
}
