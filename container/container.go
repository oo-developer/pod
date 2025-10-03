package container

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/oo-developer/pod/common"
)

//go:embed default.pod
var defaultPod string

type Config struct {
	Port       int  `json:"port"`
	FirstStart bool `json:"firstStart"`
}

func Init(config common.ConfigService, system common.SystemService) common.ContainerService {
	return &container{
		config: config,
		system: system,
	}
}

type container struct {
	config common.ConfigService
	system common.SystemService
}

func (c *container) Remove() {
	wd, _ := os.Getwd()
	c.system.Execute("podman", "rm", "-f", c.config.PodImageName(wd))
	fmt.Println("[OK] Podman image removed")
}

func (c *container) Stop() {
	wd, _ := os.Getwd()
	c.system.Execute("podman", "stop", c.config.PodImageName(wd))
	fmt.Println("[OK] Podman image stopped")
}

func (c *container) Start() {
	wd, _ := os.Getwd()
	c.system.Execute("podman", "start", c.config.PodImageName(wd))
	fmt.Println("[OK] Podman image started")
}

func (c *container) BuildDockerFile() {
	wd, _ := os.Getwd()
	podDefinition := c.GetDefaultPod()
	dockerfilePath := path.Join(c.config.PodPath(wd), "/Dockerfile")
	c.system.CheckOrCreatePath(c.config.PodPath(wd))
	c.prepareDockerfileBuild(c.config.PodPath(wd), podDefinition)
	defer c.cleanupDockerfileBuild(c.config.PodPath(wd), podDefinition)
	var docker string
	switch podDefinition.Container.Flavor {
	case "debian":
		docker = c.buildDockerfileDebian(podDefinition)
		break
	case "":
		docker = c.buildDockerfileDebian(podDefinition)
		break
	default:
		fmt.Printf("[ERROR] Unknown container flavor '%s'\n", podDefinition.Container.Flavor)
		os.Exit(1)
	}
	err := os.WriteFile(dockerfilePath, []byte(docker), 0660)
	if err != nil {
		fmt.Printf("[ERROR] Writing Dockerfile '%s':  %v\n", dockerfilePath, err)
		os.Exit(1)
	}
	fmt.Printf("[OK] Dockerfile written to '%s'\n", dockerfilePath)
	c.buildImage()
}

func (c *container) buildImage() {
	fmt.Println("[OK] Building podman image...")
	wd, _ := os.Getwd()
	podPath := c.config.PodPath(wd)
	imageName := c.config.PodImageName(wd)
	err := os.Chdir(podPath)
	defer os.Chdir(wd)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	c.system.Execute("podman", "build", "-t", imageName, ".")
	fmt.Printf("[OK] Podman image '%s' built\n", imageName)
}

func (c *container) buildDockerfileDebian(podDefinition *common.PodDefinition) string {
	docker := fmt.Sprintf("FROM %s", podDefinition.Container.Image)
	docker += "\n"
	docker += "RUN apt update && apt install -y \\\n"
	docker += "wget curl git vim lsb-release openssh-server systemd sudo "
	for _, pkg := range podDefinition.Packages {
		docker += fmt.Sprintf(" %s", pkg)
	}
	docker += "\n\n"
	docker += "RUN userdel ubuntu 2>/dev/null || true\n"
	docker += fmt.Sprintf("RUN useradd --uid %d --create-home --shell /bin/bash --no-log-init %s\n", c.system.Uid(), c.system.User())
	docker += fmt.Sprintf("RUN adduser %s sudo\n", c.system.User())
	docker += fmt.Sprintf("RUN echo '%s ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers\n", c.system.User())
	docker += fmt.Sprintf("RUN mkdir -p /home/%s/.ssh\n", c.system.User())
	docker += fmt.Sprintf("COPY id_rsa /home/%s/.ssh/id_rsa\n", c.system.User())
	docker += fmt.Sprintf("COPY id_rsa.pub /home/%s/.ssh/id_rsa.pub\n", c.system.User())
	docker += fmt.Sprintf("RUN chmod 600 /home/%s/.ssh/id_rsa\n", c.system.User())
	docker += fmt.Sprintf("RUN echo '%s' >> /home/%s/.ssh/authorized_keys\n", podDefinition.Ssh.AuthorizedKey, c.system.User())
	docker += fmt.Sprintf("RUN echo 'eval \"$(ssh-agent -s)\"' >> /home/%s/.bashrc\n", c.system.User())
	docker += fmt.Sprintf("RUN echo 'ssh-add' >> /home/%s/.bashrc\n", c.system.User())
	docker += fmt.Sprintf("RUN chown -R %s:%s /home/%s/.ssh\n", c.system.User(), c.system.User(), c.system.User())
	for _, recipeName := range podDefinition.Recipes {
		if recipe, ok := c.getUserRecipe(recipeName, "root"); ok {
			docker += recipe
		}
	}
	docker += fmt.Sprintf("USER %s\n", c.system.User())
	for _, recipeName := range podDefinition.Recipes {
		if recipe, ok := c.getUserRecipe(recipeName, "user"); ok {
			docker += recipe
		}
	}
	docker += "\n"
	docker += "USER root\n"
	docker += "ENV container=podman\n"
	docker += "CMD [\"/lib/systemd/systemd\", \"--system\"]\n"
	return docker
}

func (c *container) prepareDockerfileBuild(podPath string, podDefinition *common.PodDefinition) {
	sshPath := podDefinition.Ssh.PrivateKeyPath
	sshPrivateKeyFileSrc := path.Join(sshPath, "id_rsa")
	sshPublicKeyFileSrc := path.Join(sshPath, "id_rsa.pub")
	sshPrivateKeyFileTarget := path.Join(podPath, "id_rsa")
	sshPublicKeyFileTarget := path.Join(podPath, "id_rsa.pub")
	err := c.system.CopyFile(sshPrivateKeyFileSrc, sshPrivateKeyFileTarget)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	err = c.system.CopyFile(sshPublicKeyFileSrc, sshPublicKeyFileTarget)
	if err != nil {
		fmt.Printf("[ERROR] No private key file found\n", err)
		os.Exit(1)
	}
	err = os.Chmod(sshPrivateKeyFileTarget, 0600)
	err = c.system.CopyFile(sshPublicKeyFileSrc, sshPublicKeyFileTarget)
	if err != nil {
		fmt.Printf("[ERROR] No public key file found\n", err)
		os.Exit(1)
	}
}

func (c *container) cleanupDockerfileBuild(podPath string, podDefinition *common.PodDefinition) {
	sshPrivateKeyFileTarget := path.Join(podPath, "id_rsa")
	sshPublicKeyFileTarget := path.Join(podPath, "id_rsa.pub")
	err := os.Remove(sshPrivateKeyFileTarget)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	err = os.Remove(sshPublicKeyFileTarget)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
}

func (c *container) WriteDefaultPod() {
	wd, _ := os.Getwd()
	fileName := path.Join(wd, "/default.pod")
	if c.system.PathExists(fileName) {
		fmt.Printf("[ERROR] Pod file '%s' already exists\n", fileName)
		os.Exit(1)
	}
	data := strings.Replace(defaultPod, "%authorizedKey", c.config.AuthorizedKey(), 1)
	data = strings.Replace(data, "%privateKeyPath", path.Join(c.system.HomeDir(), ".ssh"), 1)
	err := os.WriteFile(fileName, []byte(data), 0640)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[OK] Default pod written to '%s'\n", fileName)
}

func (c *container) GetDefaultPod() *common.PodDefinition {
	wd, _ := os.Getwd()
	fileName := path.Join(wd, "/default.pod")
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	podDef := &common.PodDefinition{}
	err = json.Unmarshal(data, &podDef)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	if podDef.Container.Name == "" {
		podDef.Container.Name = "pod"
	}
	if podDef.Container.Mount == "" {
		podDef.Container.Mount = "data"
	}
	return podDef
}

func (c *container) RunContainer() {
	fmt.Println("[OK] Running container...")
	wd, _ := os.Getwd()
	defaultPod := c.GetDefaultPod()
	imageName := c.config.PodImageName(wd)
	config := c.loadContainerConfig()
	args := []string{"run"}
	args = append(args, "--name")
	args = append(args, imageName)
	args = append(args, "-p")
	args = append(args, fmt.Sprintf("%d:22", config.Port))
	args = append(args, "-d")
	args = append(args, "--privileged")
	args = append(args, "--systemd=always")
	args = append(args, "--volume")
	args = append(args, fmt.Sprintf("%s:/home/%s/%s:U", wd, c.system.User(), defaultPod.Container.Mount))
	args = append(args, "--env")
	args = append(args, fmt.Sprintf("DISPLAY=%s", os.Getenv("DISPLAY")))
	args = append(args, "--tmpfs")
	args = append(args, "/tmp")
	args = append(args, "--tmpfs")
	args = append(args, "/run")
	args = append(args, "--tmpfs")
	args = append(args, "/run/lock")
	args = append(args, "--volume")
	args = append(args, "/sys/fs/cgroup:/sys/fs/cgroup:ro")
	args = append(args, "--hostname")
	args = append(args, defaultPod.Container.Name)
	args = append(args, imageName)
	c.system.Execute("podman", args...)
	c.system.Execute("podman", "exec", "-it", imageName, "sudo", "systemctl", "start", "ssh")
	c.system.Execute("podman", "exec", "-it", imageName, "sudo", "chmod", "0777", fmt.Sprintf("/home/%s/project", c.system.User()))
	config.FirstStart = true
	c.saveContainerConfig(config)
	fmt.Printf("[OK] Container '%s' is running\n", imageName)
}

func (c *container) loadContainerConfig() *Config {
	config := &Config{}
	wd, _ := os.Getwd()
	podPath := c.config.PodPath(wd)
	configFile := path.Join(podPath, "config.json")
	if !c.system.PathExists(configFile) {
		config.Port = c.system.FreePort()
		config.FirstStart = false
		c.saveContainerConfig(config)
		return config
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	return config
}

func (c *container) saveContainerConfig(config *Config) {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	wd, _ := os.Getwd()
	podPath := c.config.PodPath(wd)
	configFile := path.Join(podPath, "config.json")
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(configFile, jsonData, 0660)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
}

func (c *container) getUserRecipe(recipeName, user string) (string, bool) {
	items, _ := os.ReadDir(c.config.RecipesPath())
	for _, item := range items {
		if !item.IsDir() && strings.HasSuffix(item.Name(), ".rcp") {
			name := strings.TrimSuffix(item.Name(), ".rcp")
			shortName := strings.TrimPrefix(name, "root_")
			shortName = strings.TrimPrefix(shortName, "user_")
			if shortName == recipeName && strings.HasPrefix(name, user) {
				data, err := os.ReadFile(path.Join(c.config.RecipesPath(), item.Name()))
				if err != nil {
					fmt.Printf("[WARNING] Loading recipe '%s': %v\n", recipeName, err)
					os.Exit(1)
				}
				return string(data), true
			}
		}
	}
	fmt.Printf("[WARNING] Recipe '%s' not found\n", recipeName)
	return "", false
}
