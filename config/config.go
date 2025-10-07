package config

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/oo-developer/pod/common"
)

func Init(system common.SystemService) common.ConfigService {
	c := &config{
		system: system,
	}
	c.load()
	return c
}

type config struct {
	system         common.SystemService
	LibraryGitRepo string `json:"libraryGitRepository"`
}

func (c *config) load() {
	exists := c.system.CheckOrCreatePath(c.BasePath())
	if !exists {
		c.LibraryGitRepo = "git@github.com:oo-developer/pod-lib.git"
		c.save()
	}
	c.system.CheckOrCreatePath(c.PodsPath())
	data, err := os.ReadFile(c.ConfigPath())
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("[OK] Config loaded from '%s'\n", c.ConfigPath())
	if !exists {
		wd, _ := os.Getwd()
		defer os.Chdir(wd)
		err = os.Chdir(c.BasePath())
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			os.Exit(1)
		}
		c.system.Execute("git", "clone", c.LibraryGitRepository())
		fmt.Printf("[OK] Repositotry '%s' cloned!\n", c.LibraryGitRepository())
		err := os.Rename(path.Join(c.BasePath(), "pod-lib"), c.LibraryPath())
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			os.Exit(1)
		}
	}
}

func (c *config) save() {
	jsonData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(c.ConfigPath(), jsonData, 0640)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[OK] ConfigService saved!")
}

func (c *config) ConfigPath() string {
	return path.Join(c.BasePath(), "config.json")
}

func (c *config) LibraryGitRepository() string {
	return "git@github.com:oo-developer/pod-lib.git"
}

func (c *config) BasePath() string {
	return path.Join(c.system.HomeDir(), ".pod")
}

func (c *config) LibraryPath() string {
	return path.Join(c.BasePath(), "library")
}

func (c *config) RecipesPath(flavor string) string {
	return path.Join(c.LibraryPath(), "recipes", flavor, c.system.Architecture())
}

func (c *config) ListRecipes(flavor string) {
	items, _ := os.ReadDir(c.RecipesPath(flavor))
	for _, item := range items {
		if !item.IsDir() && strings.HasSuffix(item.Name(), ".rcp") {
			name := strings.TrimSuffix(item.Name(), ".rcp")
			name = strings.TrimPrefix(name, "root_")
			name = strings.TrimPrefix(name, "user_")
			if strings.HasPrefix(item.Name(), "root_") {
				fmt.Printf("[root] %s\n", name)
			} else if strings.HasPrefix(item.Name(), "user_") {
				fmt.Printf("[user] %s\n", name)
			}
		}
	}
}

func (c *config) PodsPath() string {
	return path.Join(c.BasePath(), "pods")
}

func (c *config) PodPath(folder string) string {
	return path.Join(c.PodsPath(), c.PodImageName(folder))
}

func (c *config) PodImageName(folder string) string {
	hasher := sha1.New()
	hasher.Write([]byte(folder))
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func (c *config) AuthorizedKey() string {
	path := path.Join(c.system.HomeDir(), ".ssh", "id_rsa.pub")
	if c.system.PathExists(path) {
		data, _ := os.ReadFile(path)
		return strings.TrimSuffix(string(data), "\n")
	} else {
		fmt.Printf("[ERROR] No public key found: %s\n", path)
		os.Exit(1)
		return ""
	}
}

func (c *config) PrivateKeyPath() string {
	return path.Join(c.BasePath(), "id_rsa")
}
