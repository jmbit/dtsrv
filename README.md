# Project dtsrv (DockerTerminalSeRVer)

Attempt to create a terminal server style application for Docker containers, mainly intended to be used with the webtop-Containers by linuxserver.io, but should be usable with anything that listens on a HTTP endpoint.

## Getting Started
To get up and running with a (development) install of dtsrv, you need go, templ, air and docker installed on the system.
you also need to create a `.env` file in the project directory like this:

```env
PORT=8080
APP_ENV=local
SESSION_KEY=randomkey
IMAGE_NAME=nameofdockerimage
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite (as soon as it actually exists)
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
