# vvp2 CLI

> **⚠️ EXPERIMENTAL PROJECT**: This is an experimental CLI tool under active development. Features may be incomplete, APIs may change, and stability is not guaranteed. Use at your own risk and avoid using in production environments without thorough testing.
> 
> **⚠️ UNOFFICIAL PROJECT**: This project is **not official** and is **not supported by Ververica**. It is a community-driven tool created independently. For official support and tooling, please refer to the [Ververica Platform documentation](https://docs.ververica.com/).

A command-line interface tool for interacting with the Ververica Platform (VVP) API. Built with Go, using Cobra for CLI framework and Viper for configuration management.

## ⚠️ API Status

**Latest Update**: API endpoints have been validated against the VVP OpenAPI specification. Some features may have limited functionality:

- ✅ **Session Clusters**: Fully functional via Application Manager API (`/api/v1/namespaces/{ns}/sessionclusters`)
 

## Features

- **Session Cluster Management**: Manage SQL session clusters (list, get, create, update, delete)
- **Job Management**: List and view Flink jobs running in the platform
- **TLS Options**: Support for insecure mode to skip TLS verification

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/mcolomerc/vvp2-cli/main/install.sh | bash
```

### Manual Binary Download

Download the appropriate binary for your platform from the [releases page](https://github.com/mcolomerc/vvp2-cli/releases):

- **Linux (AMD64)**: `vvp2-linux-amd64.tar.gz`
- **Linux (ARM64)**: `vvp2-linux-arm64.tar.gz`
- **macOS (Intel)**: `vvp2-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `vvp2-darwin-arm64.tar.gz`
- **Windows (AMD64)**: `vvp2-windows-amd64.zip`

Then extract and install:

```bash
# Example for Linux AMD64
tar -xzf vvp2-linux-amd64.tar.gz
sudo mv vvp2 /usr/local/bin/
chmod +x /usr/local/bin/vvp2
```

### Build from Source

```bash
git clone https://github.com/mcolomerc/vvp2-cli.git
cd vvp2-cli
make build
# Or install to GOPATH/bin
make install
```

### Verify Installation

```bash
vvp2 version
```

## Quick Start

After installation, initialize your configuration:

```bash
# Run the interactive configuration wizard
./vvp2 config init

# Or set via environment variables
export VVP_API_URL="http://vvp.localhost"
export VVP_API_TOKEN="your-api-token"

# Start using the CLI
./vvp2 namespace list
```

## Configuration

The CLI can be configured in three ways (in order of precedence):

1. **Command-line flags**
2. **Environment variables** (prefixed with `VVP_`)
3. **Configuration file** (`~/.vvp2/config.yaml` or specified with `--config`)

### Interactive Configuration Setup

The easiest way to get started is to use the interactive configuration wizard:

```bash
# Initialize configuration interactively
vvp2 config init

# View your configuration
vvp2 config show

# Show configuration file path
vvp2 config path

# Overwrite existing configuration
vvp2 config init --force
```

The wizard will prompt you for:
- Ververica Platform API URL
- API Token (optional)
- TLS certificate verification settings
- Default namespace
- Default output format

### Configuration File

The configuration is stored at `~/.vvp2/config.yaml`. You can also manually create or edit this file:

```yaml
api:
  url: "http://vvp.localhost"
  token: "your-api-token"
  insecure: false

default:
  namespace: "default"

output:
  format: "table"  # table, json, or yaml
```

### Environment Variables

```bash
export VVP_API_URL="http://vvp.localhost"
export VVP_API_TOKEN="your-api-token"
export VVP_API_INSECURE="false"
export VVP_DEFAULT_NAMESPACE="default"
export VVP_OUTPUT_FORMAT="table"
```

### Command-line Flags

```bash
vvp2 --api-url http://vvp.localhost --api-token your-token --namespace default
```

## Usage

### Global Flags

- `--api-url`: Ververica Platform API URL
- `--api-token`: API authentication token
- `--namespace`: Default namespace
- `--insecure`: Skip TLS certificate verification
- `--output, -o`: Output format (table, json, yaml)
- `--config`: Config file path (default: `$HOME/.vvp2/config.yaml`)

### Configuration Commands

```bash
# Initialize configuration interactively
vvp2 config init

# Show current configuration
vvp2 config show

# Show configuration file path
vvp2 config path

# Reinitialize (overwrite) configuration
vvp2 config init --force
```

### Namespace Commands

```bash
# List all namespaces
vvp2 namespace list

# Get a specific namespace
vvp2 namespace get my-namespace

# Create a namespace from file
vvp2 namespace create -f namespace.yaml

# Update a namespace
vvp2 namespace update my-namespace -f namespace.yaml

# Delete a namespace
vvp2 namespace delete my-namespace
```

### Deployment Commands

Note: If you configured a default namespace (via `vvp2 config init` or `~/.vvp2/config.yaml`), you can omit `-n/--namespace`.

```bash
# List deployments in a namespace
vvp2 deployment list -n my-namespace

# Get a specific deployment
vvp2 deployment get my-deployment -n my-namespace

# Create a deployment from file
vvp2 deployment create -n my-namespace -f deployment.yaml

# Update a deployment
vvp2 deployment update my-deployment -n my-namespace -f deployment.yaml

# Delete a deployment
vvp2 deployment delete my-deployment -n my-namespace

# Start a deployment
vvp2 deployment start my-deployment -n my-namespace

# Stop a deployment
vvp2 deployment stop my-deployment -n my-namespace

# Suspend a deployment
vvp2 deployment suspend my-deployment -n my-namespace
```

### Deployment Target Commands

Note: If you configured a default namespace (via `vvp2 config init` or `~/.vvp2/config.yaml`), you can omit `-n/--namespace`.

```bash
# List deployment targets in a namespace
vvp2 deployment-target list -n my-namespace
# Or use alias
vvp2 dt list -n my-namespace

# Get a specific deployment target
vvp2 dt get my-target -n my-namespace

# Create a deployment target from file
vvp2 dt create -n my-namespace -f deploymenttarget.yaml

# Update a deployment target
vvp2 dt update my-target -n my-namespace -f deploymenttarget.yaml

# Delete a deployment target
vvp2 dt delete my-target -n my-namespace
```

### Deployment Defaults Commands

Note: Defaults are namespaced. If you configured a default namespace, you can omit `-n/--namespace`.

```bash
# Get deployment defaults for a namespace
vvp2 deployment-defaults get -n my-namespace

# Replace deployment defaults from file (YAML/JSON)
vvp2 deployment-defaults replace -n my-namespace -f defaults.yaml

# Update deployment defaults via PATCH with a SecretValue (advanced)
vvp2 deployment-defaults update -n my-namespace -f secretvalue.yaml
```

### Session Cluster Commands

Session clusters provide interactive SQL environments for executing Flink SQL queries. If you configured a default namespace, you can omit `-n/--namespace`.

```bash
# List session clusters in a namespace
vvp2 sessioncluster list -n my-namespace
# Or use aliases
vvp2 sc list -n my-namespace
vvp2 session-cluster list -n my-namespace

# Get a specific session cluster
vvp2 sc get my-sql-session -n my-namespace

# Create a session cluster from file
vvp2 sc create -n my-namespace -f sessioncluster.yaml

# Update a session cluster
vvp2 sc update my-sql-session -n my-namespace -f sessioncluster.yaml

# Delete a session cluster
vvp2 sc delete my-sql-session -n my-namespace
```

### Job Commands

Jobs represent running Flink applications in the platform. Use these commands to view job status and details.

```bash
# List all jobs in a namespace
vvp2 job list -n my-namespace
# Or use alias
vvp2 jobs list -n my-namespace

# Get details of a specific job by ID
vvp2 job get <job-id> -n my-namespace

# List jobs with JSON output
vvp2 job list -n my-namespace -o json

# Get job details in YAML format
vvp2 job get <job-id> -n my-namespace -o yaml
```

**Note**: Job IDs are UUIDs (e.g., `e998a415-1d6e-4a97-bd64-590f20b605e7`). You can get the job ID from the `job list` command output.

### Resource Usage Report Command

If your Ververica Platform instance has resource usage tracking enabled, you can generate a platform-wide resource usage report:

```bash
vvp2 usage report
```
- By default, this returns usage for the last 7 days.
- You can specify a custom date range with `--from` and `--to` (YYYY-MM-DD format):

```bash
vvp2 usage report --from "2025-11-01" --to "2025-11-18"
```
- Dates are in LocalDate format (`YYYY-MM-DD`), interpreted as 00:00 UTC.
- `--from` is inclusive, `--to` is exclusive.
- The endpoint is `/api/v1/status/resourceusage` (platform-wide, not namespaced).
- Output format can be set with `-o json`, `-o yaml`, or default table view.
- If the endpoint is not enabled, you'll see a clear error message.

See [VVP documentation](https://docs.ververica.com/vvp/platform-operations/advanced-configurations/resource-usage-tracking/?highlight=Usage#generating-a-resource-usage-report) for details on enabling this feature.

## Example Files

### Namespace Example

```yaml
metadata:
  name: my-namespace
  labels:
    env: production
spec:
  roleBindings:
    - role: owner
      members:
        - user@example.com
```

  ### Deployment Target Example

  ```yaml
  metadata:
    name: kubernetes-target
    namespace: my-namespace
    labels:
      environment: production
  spec:
    kubernetes:
      namespace: "flink-jobs"
  ```

### Session Cluster Example

```yaml
apiVersion: v1
kind: SessionCluster
metadata:
  name: my-sql-session
  namespace: default
  labels:
    environment: dev
    team: data-platform
spec:
  deploymentTargetName: vvp-jobs-flink
  state: RUNNING
  flinkImageRegistry: docker.io
  flinkImageRepository: confluentinc/cp-flink
  flinkImageTag: 1.19.1-cp1
  flinkVersion: "1.19"
  numberOfTaskManagers: 2
  resources:
    jobmanager:
      cpu: 1.0
      memory: 2048m
    taskmanager:
      cpu: 2.0
      memory: 4096m
```

### Deployment Example

```yaml
metadata:
  name: my-deployment
  namespace: my-namespace
  labels:
    app: flink-job
spec:
  state: RUNNING
  upgradeStrategy:
    kind: STATEFUL
  template:
    spec:
      artifact:
        kind: JAR
        jarUri: "s3://bucket/my-job.jar"
        mainClass: "com.example.MyJob"
      parallelism: 2
      flinkVersion: "1.17"
      flinkConfiguration:
        taskmanager.numberOfTaskSlots: "2"
      resources:
        jobmanager:
          cpu: "1"
          memory: "1G"
        taskmanager:
          cpu: "2"
          memory: "2G"
```

## Output Formats

### Table (default)
```bash
vvp-cli deployment list -n my-namespace
```

### JSON
```bash
vvp-cli deployment get my-deployment -n my-namespace -o json
```

### YAML
```bash
vvp-cli deployment get my-deployment -n my-namespace -o yaml
```

## Development

### Project Structure

```
.
├── main.go                 # Application entry point
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and configuration
│   ├── config.go         # Configuration commands
│   ├── deployment.go      # Deployment commands
│   ├── deploymenttarget.go # Deployment target commands
│   ├── namespace.go       # Namespace commands
│   └── session.go         # Session commands
└── pkg/                    # Reusable packages
    ├── api/               # API client and models
    │   ├── client.go      # HTTP client setup
    │   ├── deployment.go  # Deployment API methods
  │   ├── deploymentdefaults.go # Deployment defaults API methods
    │   ├── deploymenttarget.go # Deployment target API methods
    │   ├── namespace.go   # Namespace API methods
    │   └── session.go     # Session API methods
    └── config/            # Configuration management
        └── config.go      # Config structures and loading
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra): CLI framework
- [Viper](https://github.com/spf13/viper): Configuration management
- [Resty](https://github.com/go-resty/resty): HTTP client

### Build

```bash
go build -o vvp2
```

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
