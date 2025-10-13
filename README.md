# POD

POD is wrapper around podman to create and run containers on mainly desktop Linux.
You can run ui applications like development UIs or browsers from inside the container.

The main goal of this project is to keep your host system clean.

This project is in an early stage!

## Install

Requirements:
- go
- git
- podman
- ssh key (id_rsa and id_rsa.pub) in $HOME/.ssh/

```
git clone git@github.com:oo-developer/pod.git
cd pod
go build -o pod
sudo cp pod /usr/local/bin/
rm pod
```

## Create and run a POD

A pod ha a strong entanglement with a directory. So first create a directory.
```
mkdir nostromo
cd nostromo
```

Then initialize the POD
```
pod init
```
A new default.pod file will be created.
```
{
    "container": {
        "image": "ubuntu:24.04",
        "flavor": "debian",
        "name": "nostromo",
        "mount": "project"
    },
    "ssh": {
        "privateKeyPath": "<path to a folder with id_rsa and id_rsa.pub you want to use from inside container>",
        "authorizedKey": "<ssh public key from your %HOME/.ssh/id_rs.pub>"
    },
    "packages": [
        "x11-apps"
    ],
    "recipes": [
        "vscode",
        "vivaldi",
        "go"
    ]
}
```
- The .container.mount entry will create a folder on the host which will be mounted in the container.
- In .packages you can add additional packages you want to use
- A recipe is software which can not be installed via packages

After editing default.pod you can build the container image. This can take some time. 
```
pod build
```

The run the POD
```
pod run
```

## Access the pod

### Using the shell command

```
pod shell
```
Or outside the pod directory
```
pod shell nostromo
```

### Via ssh

First look for the ssh port of the POD.
```
pod status 
```
Or outside the pod directory
```
pod status nostromo
```

Or use 
```
pod list
```

Open a ssh session to the pod
```
ssh -p <port> -X $USER@localhost 
```