#!/bin/sh
# Local test script for building Cribl Storage Tool for multiple platforms

# Check if we have Go files to build
if [ ! "$(find . -name "*.go" | head -n 1)" ]; then
  echo "Error: No Go files found in current directory or subdirectories."
  echo "Make sure you're running this from your Go project root."
  exit 1
fi

# Identify the main package
MAIN_PACKAGE=""
for file in $(find . -name "*.go" -type f); do
  if grep -q "func main()" "$file"; then
    dir=$(dirname "$file")
    MAIN_PACKAGE=$(realpath --relative-to="$(pwd)" "$dir")
    echo "Found main package in: $MAIN_PACKAGE"
    break
  fi
done

if [ -z "$MAIN_PACKAGE" ]; then
  echo "Warning: No main() function found. Will try to build from current directory."
  MAIN_PACKAGE="."
fi

# Create a directory for the builds
mkdir -p builds

# Build for each platform (OS/architecture)
build_platform() {
  platform=$1

  # Split the platform string into OS and architecture
  GOOS=$(echo "$platform" | cut -d'/' -f1)
  GOARCH=$(echo "$platform" | cut -d'/' -f2)

  # Create platform-specific directory
  mkdir -p "builds/${GOOS}-${GOARCH}"

  # Set the output filename - same name for all platforms
  if [ "$GOOS" = "windows" ]; then
    output_name="builds/${GOOS}-${GOARCH}/cribl-storage-tool.exe"
  else
    output_name="builds/${GOOS}-${GOARCH}/cribl-storage-tool"
  fi

  # Build the binary
  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name" "$MAIN_PACKAGE"

  # Check if the build was successful
  if [ $? -eq 0 ]; then
    echo "✓ Successfully built for $GOOS/$GOARCH"
  else
    echo "✗ Failed to build for $GOOS/$GOARCH"
    exit 1
  fi
}

# Build for all platforms
build_platform "linux/amd64"
build_platform "linux/arm64"
build_platform "windows/amd64"
build_platform "darwin/amd64"
build_platform "darwin/arm64"
build_platform "freebsd/amd64"

echo "Build process completed!"
echo "All binaries are named 'cribl-storage-tool' within their platform directories:"
find builds -type f | sort