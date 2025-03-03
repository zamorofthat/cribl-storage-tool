.PHONY: all deps test build clean

# Go parameters
GO = go

all: deps test build

deps:
	$(GO) mod download
	@echo "Dependencies installed"

test:
	$(GO) test -v ./...

build: clean
	@mkdir -p build
	@echo "Building for all platforms..."

	@echo "Building for linux/amd64..."
	@mkdir -p build/linux-amd64
	@GOOS=linux GOARCH=amd64 $(GO) build -o build/linux-amd64/cribl-storage-tool .
	@echo "✓ Successfully built for linux/amd64"

	@echo "Building for linux/arm64..."
	@mkdir -p build/linux-arm64
	@GOOS=linux GOARCH=arm64 $(GO) build -o build/linux-arm64/cribl-storage-tool .
	@echo "✓ Successfully built for linux/arm64"

	@echo "Building for windows/amd64..."
	@mkdir -p build/windows-amd64
	@GOOS=windows GOARCH=amd64 $(GO) build -o build/windows-amd64/cribl-storage-tool.exe .
	@echo "✓ Successfully built for windows/amd64"

	@echo "Building for darwin/amd64..."
	@mkdir -p build/darwin-amd64
	@GOOS=darwin GOARCH=amd64 $(GO) build -o build/darwin-amd64/cribl-storage-tool .
	@echo "✓ Successfully built for darwin/amd64"

	@echo "Building for darwin/arm64..."
	@mkdir -p build/darwin-arm64
	@GOOS=darwin GOARCH=arm64 $(GO) build -o build/darwin-arm64/cribl-storage-tool .
	@echo "✓ Successfully built for darwin/arm64"

	@echo "Building for freebsd/amd64..."
	@mkdir -p build/freebsd-amd64
	@GOOS=freebsd GOARCH=amd64 $(GO) build -o build/freebsd-amd64/cribl-storage-tool .
	@echo "✓ Successfully built for freebsd/amd64"

	@echo "Build process completed!"

clean:
	@rm -rf build/
	@echo "Build directory cleaned"