WORKY        := $(CURDIR)
SCAFFOLD_DIR ?= /tmp/worky-scaffold

.PHONY: build dev clean-scaffold

build:
	go build -o bin/worky ./cmd/worky

## dev: build the CLI and run `worky init` interactively.
## Each run creates a new workshop under SCAFFOLD_DIR/<slug>.
## Override the base dir with: make dev SCAFFOLD_DIR=/my/path
dev: build
	SCAFFOLD_DIR=$(SCAFFOLD_DIR) ./scripts/dev-init.sh

## clean-scaffold: remove all generated test workshops
clean-scaffold:
	rm -rf "$(SCAFFOLD_DIR)"
	@echo "Removed $(SCAFFOLD_DIR)"
