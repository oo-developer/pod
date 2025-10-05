package container

import "fmt"

var packages = []string{
	"wget",
	"curl",
	"git",
	"vim",
	"lsb-release",
	"openssh-server",
	"systemd",
	"sudo",
	"gpg",
	"apt-transport-https",
	"ca-certificates",
	"gnupg",
	"bzip2",
	"libpci3",
	//"libasound2",
	"libatk1.0-0",
	"libcairo2",
	"libcups2",
	"libdbus-1-3",
	"libexpat1",
	"libfontconfig1",
	"libfreetype6",
	"libgbm1",
	"libgtk-3-0",
	"libnspr4",
	"libnss3",
	"libpango-1.0-0",
	"libx11-6",
	"libx11-xcb1",
	"libxcb1",
	"libxcomposite1",
	"libxcursor1",
	"libxdamage1",
	"libxext6",
	"libxfixes3",
	"libxi6",
	"libxrandr2",
	"libxrender1",
	"libxss1",
	"libxtst6",
	"libxshmfence-dev",
	"openssh-client",
	"xauth",
	"net-tools",
	"xz-utils",
}

func getDefaultPackages() string {
	var pkgs string
	for _, pkg := range packages {
		pkgs += fmt.Sprintf("    %s\\\n", pkg)
	}
	return pkgs
}
