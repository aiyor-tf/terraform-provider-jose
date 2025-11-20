#!/bin/bash
set -e

# Versions to test (matching CI matrix 1.0.* - 1.4.*)
VERSIONS=("1.0.11" "1.1.9" "1.2.9" "1.3.9" "1.4.6")
OS="linux"
ARCH="amd64"
BIN_DIR="$(pwd)/bin"
mkdir -p "$BIN_DIR"

for VERSION in "${VERSIONS[@]}"; do
    echo "------------------------------------------------"
    echo "Testing with Terraform $VERSION"
    echo "------------------------------------------------"

    TF_BIN="$BIN_DIR/terraform-$VERSION"
    
    if [ ! -f "$TF_BIN" ]; then
        echo "Downloading Terraform $VERSION..."
        curl -sSL -o "/tmp/terraform_${VERSION}_${OS}_${ARCH}.zip" "https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_${OS}_${ARCH}.zip"
        unzip -q -o "/tmp/terraform_${VERSION}_${OS}_${ARCH}.zip" -d "/tmp"
        mv "/tmp/terraform" "$TF_BIN"
        rm "/tmp/terraform_${VERSION}_${OS}_${ARCH}.zip"
    fi

    # Create a temporary directory for this run's PATH to ensure we use the correct version
    TMP_PATH_DIR=$(mktemp -d)
    ln -s "$TF_BIN" "$TMP_PATH_DIR/terraform"
    
    # Run tests with modified PATH
    # We run in a subshell to avoid polluting the parent shell's PATH/env
    (
        export PATH="$TMP_PATH_DIR:$PATH"
        export TF_ACC=1
        
        echo "Using binary: $(which terraform)"
        terraform version
        
        go test -v -cover ./internal/provider/
    )
    
    # Cleanup temp dir
    rm -rf "$TMP_PATH_DIR"
    
    echo "Passed Terraform $VERSION"
    echo ""
done

echo "All tests passed!"
