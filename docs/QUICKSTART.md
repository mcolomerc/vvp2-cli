# Quick Start Guide

This guide will help you get started with the VVP CLI tool.

## Prerequisites

- Go 1.21 or later (for building from source)
- Access to a Ververica Platform API endpoint
- API token for authentication (optional, depending on your VVP setup)

## Installation

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd vvp2-cli

# Build the binary
make build

# Or using go directly
go build -o vvp2

# Optionally install to your PATH
make install
```

## Configuration

### Quick Setup (Recommended)

The easiest way to configure vvp2 is using the interactive setup wizard:

```bash
./vvp2 config init
```

This will prompt you for:
- API URL
- API Token (optional)
- TLS settings
- Default namespace
- Output format preference

The configuration will be saved to `~/.vvp2/config.yaml`.

### Option 1: Interactive Configuration (Recommended)

```bash
# Initialize configuration
./vvp2 config init

# View your configuration
./vvp2 config show

# Check configuration file location
./vvp2 config path

# Reconfigure (overwrite existing)
./vvp2 config init --force
```

### Option 2: Manual Configuration File

Create the configuration directory and file at `~/.vvp2/config.yaml`:

```yaml
api:
  url: "http://vvp.localhost"
  token: "your-api-token"
  insecure: false

default:
  namespace: "default"

output:
  format: "table"
```

You can copy the example configuration:

```bash
# Create the config directory
mkdir -p ~/.vvp2

# Copy the example config
cp config.yaml.example ~/.vvp2/config.yaml

# Edit with your settings
nano ~/.vvp2/config.yaml
```

### Option 3: Environment Variables

```bash
export VVP_API_URL="http://vvp.localhost"
export VVP_API_TOKEN="your-api-token"
export VVP_DEFAULT_NAMESPACE="default"
```

### Option 4: Command-line Flags

```bash
./vvp2 --api-url http://vvp.localhost --api-token your-token namespace list
```

## Basic Usage Examples

### Working with Namespaces

```bash
# List all namespaces
./vvp2 namespace list

# Get details of a specific namespace
./vvp2 namespace get production

# Create a namespace from a YAML file
./vvp2 namespace create -f examples/namespace.yaml

# Delete a namespace
./vvp2 namespace delete test-namespace
```

### Working with Deployments

```bash
# List deployments in a namespace
./vvp2 deployment list -n production

# Get deployment details
./vvp2 deployment get my-job -n production

# Create a deployment
./vvp2 deployment create -n production -f examples/deployment.yaml

# Start a deployment
./vvp2 deployment start my-job -n production

# Stop a deployment
./vvp2 deployment stop my-job -n production

# Suspend a deployment (with savepoint)
./vvp2 deployment suspend my-job -n production

# Delete a deployment
./vvp2 deployment delete my-job -n production
```

### Working with Deployment Targets

```bash
# List deployment targets in a namespace
./vvp2 deployment-target list -n production
# Or use the shorter alias
./vvp2 dt list -n production

# Get deployment target details
./vvp2 dt get kubernetes-target -n production

# Create a deployment target
./vvp2 dt create -n production -f examples/deploymenttarget.yaml

# Delete a deployment target
./vvp2 dt delete kubernetes-target -n production
```

### Working with Sessions

```bash
# List sessions in a namespace
./vvp2 session list -n production

# Get session details
./vvp2 session get sql-session -n production

# Create a session
./vvp2 session create -n production -f examples/session.yaml

# Delete a session
./vvp2 session delete sql-session -n production
```

### Output Formats

You can specify the output format using the `-o` or `--output` flag:

```bash
# Table format (default)
./vvp2 namespace list

# JSON format
./vvp2 namespace list -o json

# YAML format
./vvp2 namespace get production -o yaml
```

## Common Workflows

### Deploying a New Flink Job

1. Prepare your deployment YAML file (see `examples/deployment.yaml`)
2. Create the deployment:
   ```bash
   ./vvp2 deployment create -n production -f my-deployment.yaml
   ```
3. Check the status:
   ```bash
   ./vvp2 deployment get my-job -n production
   ```

### Updating a Deployment

1. Modify your deployment YAML file
2. Update the deployment:
   ```bash
   ./vvp2 deployment update my-job -n production -f my-deployment.yaml
   ```

### Managing Deployment State

```bash
# Start a stopped deployment
./vvp2 deployment start my-job -n production

# Stop a running deployment
./vvp2 deployment stop my-job -n production

# Suspend with savepoint
./vvp2 deployment suspend my-job -n production
```

## Troubleshooting

### Connection Issues

If you're having connection issues:

1. Check your API URL is correct
2. For self-signed certificates, use the `--insecure` flag:
   ```bash
   ./vvp2 --insecure namespace list
   ```
3. Verify your API token is valid

### Authentication Issues

If you get authentication errors:

1. Ensure your API token is set correctly
2. Check token permissions in VVP
3. Verify the token hasn't expired

### Debug Mode

For more verbose output, you can check the HTTP requests and responses by setting debug mode in the code or checking the API responses.

## Next Steps

- Read the full [README.md](../README.md) for detailed documentation
- Check the [examples](../examples/) directory for sample resource definitions
- Explore the Ververica Platform API documentation at your VVP instance: `http://your-vvp-instance/swagger`

## Getting Help

```bash
# General help
./vvp2 --help

# Help for specific commands
./vvp2 deployment --help
./vvp2 namespace --help
./vvp2 session --help

# Help for subcommands
./vvp2 deployment create --help
```
