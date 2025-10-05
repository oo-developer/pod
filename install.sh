#!/usr/bin/bash

go build -o build/pod
cd build || exit
sudo cp pod /usr/local/bin/
