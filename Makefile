.PHONY: all deps test build clean

# Go parameters
GOBIN ?= $(GOPATH)/bin
GO = go
PLATFORMS = linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64 freebsd/amd64

all: deps test build

deps:
	$(GO) mod download
	@echo "Dependencies installed"

test:
	$(GO) test -v ./...

build: clean
	@mkdir -p build
	@echo "Building for all platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		mkdir -p build/$$GOOS-$$GOARCH; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		if [ "$$GOOS" = "windows" ]; then \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build -o build/$$GOOS-$$GOARCH/cribl-storage-tool.exe ./cmd/cribl-storage-tool; \
		else \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build -o build/$$GOOS-$$GOARCH/cribl-storage-tool ./cmd/cribl-storage-tool; \
		fi; \
		if [ $$? -eq 0 ]; then \
			echo "✓ Successfully built for $$GOOS/$$GOARCH"; \
		else \
			echo "✗ Failed to build for $$GOOS/$$GOARCH"; \
			exit 1; \
		fi; \
	done

clean:
	@rm -rf build/
	@echo "Build directory cleaned"