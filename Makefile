# Simple Makefile for a Go project

# Build the application
all: deps build


deps:
	@go mod tidy

build:
	@echo "Building..."
	@templ generate
	@go build -o dtsrv cmd/run/main.go

# Run the application
run:
	@go run cmd/run/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f dtsrv
	@rm -f cmd/web/*_templ.go

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

install:
	@install -d /usr/local/bin; \
	useradd -r -d /var/local/lib/dtsrv dtsrv; \
	install -m 755 dtsrv /usr/local/bin/; \
	install -d -o root -g root /etc/default; \
	install -d -o root -g root example.env /etc/default/dtsrv; \
	install -d -o root -g root dtsrv.service /etc/systemd/system/dtsrv.service; \
	chown dtsrv:dtsrv /var/local/lib/dtsrv; \
  chmod 700 /var/local/lib/dtsrv; \
	usermod -aG docker dtsrv; \

	

.PHONY: all build run test clean install
