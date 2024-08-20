# Define variables
PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/man/man1

# Default target
all: build

# Build target
build:
	@echo "Building the project..."
	go build -o conman

# Install target
install: all
	@echo "Installing..."
	install -Dm755 conman $(BINDIR)/conman

# Uninstall target
uninstall:
	@echo "Uninstalling..."
	rm -f $(BINDIR)/conman

# Clean target
clean:
	@echo "Cleaning up..."
	rm -f conman

# PHONY targets
.PHONY: all build install uninstall clean

