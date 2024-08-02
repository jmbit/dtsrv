# Installing dtsrv

## Prerequisites
dtsrv requires a linux-based OS and Docker (or compatible) to work. To build you also need go 1.22 or later installed, along with templ (the templating engine).

Installing Docker: [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/)
Installing go: [https://go.dev/doc/install](https://go.dev/doc/install)
Installing templ: `go install github.com/a-h/templ/cmd/templ@latest`

## Installing dtsrv
First, build the project using `make`. If this was successful, you can then proceed to run `sudo make install` to install dtsrv on the system. 
in the next step, you need to modify `/etc/default/dtsrv` to configure dtsrv for your environment

