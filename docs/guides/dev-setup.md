# Developer Setup Guide

This guide explains how to build and test the `terraform-provider-jose` locally without publishing to the Terraform Registry.

## Prerequisites

-   Go 1.21+
-   Terraform 1.0+

## 1. Build and Install

The `GNUmakefile` includes an `install` target that compiles the provider and installs it to your `$GOBIN` (usually `$GOPATH/bin` or `~/go/bin`).

```bash
make install
```

Verify the installation:
```bash
ls $(go env GOPATH)/bin/terraform-provider-jose
```

## 2. Configure Terraform CLI

To tell Terraform to use your locally built provider instead of fetching it from the registry, you need to configure `dev_overrides` in your CLI configuration file (`~/.terraformrc` on Linux/macOS).

Add the following block to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "aiyor-tf/jose" = "/home/tliang/go/bin" # Replace with your actual GOBIN path
  }

  # For other providers, use direct installation
  direct {}
}
```

**Note:** When `dev_overrides` is active, `terraform init` will show a warning, which is expected.

## 3. Run Local Test

Navigate to the local test example directory:

```bash
cd examples/local-test
```

Initialize Terraform (this will skip provider installation for `jose` due to overrides):

```bash
terraform init
```

Apply the configuration:

```bash
terraform apply
```

You should see the resources being created and the outputs displaying the verified JWT claims and JWE.

## 4. Running Tests

To run the unit and acceptance tests:

```bash
make testacc
```
