#!/bin/bash

# Build script for multiple platforms
set -e

VERSION=${1:-"dev"}
BUILD_DIR="dist"

# Clean build directory
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for multiple platforms
platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="netmon-${GOOS}-${GOARCH}"
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for ${GOOS}/${GOARCH}..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${output_name} .
    
    # Create archive
    if [ $GOOS = "windows" ]; then
        zip -j ${BUILD_DIR}/${output_name}.zip ${BUILD_DIR}/${output_name}
    else
        tar -czf ${BUILD_DIR}/${output_name}.tar.gz -C ${BUILD_DIR} ${output_name}
    fi
    
    echo "Built ${BUILD_DIR}/${output_name}"
done

echo "Build complete! Files are in ${BUILD_DIR}/"

