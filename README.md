# sctl - Saved CLI

Command-line interface for managing Saved workspaces, agents, jobs, backups, and billing.

## Features

- 🔐 **OIDC Device Flow Authentication** with dynamic configuration
- 🏢 **Workspace Management** - Create, list, update, delete workspaces
- 🤖 **Agent Management** - Manage backup agents
- 💼 **Job Management** - Worker, agent, and manual jobs
- 💾 **Backup Operations** - Request, download, and manage backups
- 💳 **Billing & Credits** - View usage, invoices, and credit balance

## Installation

```bash
go build -o sctl ./cmd/main.go
```

## Quick Start

### 1. Initialize Configuration

```bash
./sctl config init
```

This creates `~/.sctl/config.yaml`

### 2. Configure Auth (Optional)

`sctl` will automatically fetch authentication configuration from the backend if not provided. Use manual configuration only if you need to override the defaults.

Set credentials in config:

```bash
./sctl config set auth_issuer https://your-oidc-provider.com/
./sctl config set auth_client_id your-client-id
./sctl config set auth_audience https://api.saved.sh
```

Or use environment variables:
```bash
export SAVED_AUTH_ISSUER=https://your-oidc-provider.com/
export SAVED_AUTH_CLIENT_ID=your-client-id
export SAVED_AUTH_AUDIENCE=https://api.saved.sh
```

### 3. Login

```bash
./sctl auth login
```

This will:
1. Automatically fetch auth config from server (if needed)
2. Display a device code and URL
2. Open your browser to authenticate
3. Save JWT token to `~/.sctl/config.yaml`

### 4. Set Default Workspace (Optional)

```bash
./sctl config set workspace_id <your-workspace-id>
```

### 5. Use Commands

```bash
# List workspaces
./sctl workspace list

# List agents (uses default workspace from config)
./sctl agent list

# Or specify workspace explicitly
./sctl agent list -w <workspace-id>

# Create a job
./sctl job create-worker -w <workspace-id> -n "Daily Backup"

# List jobs
./sctl job list -w <workspace-id>

# View billing
./sctl billing info -w <workspace-id>
./sctl billing credits balance -w <workspace-id>
```

## Commands

### Auth Commands

```bash
sctl auth login          # Login with OIDC device flow
sctl auth logout         # Remove stored credentials
sctl auth status         # Check authentication status
```

### Config Commands

```bash
sctl config init         # Initialize configuration file
sctl config set <key> <value>  # Set configuration value
sctl config get <key>    # Get configuration value
sctl config list         # List all configuration
```

### Workspace Commands

```bash
sctl workspace list                    # List all workspaces
sctl workspace get <workspace-id>      # Get workspace details
sctl workspace create                  # Create workspace
sctl workspace update <workspace-id>   # Update workspace
sctl workspace delete <workspace-id>   # Delete workspace
```

### Agent Commands

```bash
sctl agent list -w <workspace-id>              # List agents
sctl agent get <agent-id> -w <workspace-id>    # Get agent details
sctl agent create -w <workspace-id> -n "Name"  # Create agent
sctl agent update <agent-id> -w <workspace-id> # Update agent
sctl agent delete <agent-id> -w <workspace-id> # Delete agent
sctl agent reset-credentials <agent-id>        # Reset credentials
```

### Job Commands

```bash
sctl job list -w <workspace-id>                # List jobs
sctl job get <job-id> -w <workspace-id>        # Get job details
sctl job create-worker -w <workspace-id>       # Create worker job
sctl job create-agent -w <workspace-id>        # Create agent job
sctl job create-manual -w <workspace-id>       # Create manual job
sctl job update <job-id> -w <workspace-id>     # Update job
sctl job delete <job-id> -w <workspace-id>     # Delete job
sctl job trigger <job-id> -w <workspace-id>    # Trigger job
```

### Backup Commands

```bash
sctl backup list <job-id> -w <workspace-id>              # List backups
sctl backup get <job-id> <backup-id> -w <workspace-id>   # Get backup
sctl backup delete <job-id> <backup-id>                  # Delete backup
sctl backup request <job-id> -w <workspace-id>           # Request backup
sctl backup download <job-id> <backup-id>                # Get download URL
```

### Billing Commands

```bash
sctl billing info -w <workspace-id>             # Get billing info
sctl billing usage -w <workspace-id>            # Get usage history
sctl billing invoices -w <workspace-id>         # List invoices
sctl billing credits balance -w <workspace-id>  # Get credit balance
sctl billing credits transactions               # List transactions
```

## Configuration

Config file location: `~/.sctl/config.yaml`

```yaml
# Auth0 JWT token (set via 'sctl auth login')
api_key: eyJhbGc...

# Refresh token for token renewal
refresh_token: v1.M...

# Default workspace ID (optional)
workspace_id: ws_123abc

# API server URL
server_url: https://api.saved.sh

# Auth configuration (automatically fetched/saved)
auth_issuer: https://saved.us.auth0.com/
auth_client_id: your-client-id
auth_audience: https://api.saved.sh
```

## Environment Variables

All config values can be set via environment variables with `SAVED_` prefix:

```bash
SAVED_API_KEY=<token>
SAVED_WORKSPACE_ID=<workspace-id>
SAVED_SERVER_URL=https://api.saved.sh
SAVED_AUTH_ISSUER=<issuer-url>
SAVED_AUTH_CLIENT_ID=<client-id>
SAVED_AUTH_AUDIENCE=<audience>
```

## Architecture

### Context-Based Dependency Injection

All commands receive an authenticated `CLIContext` with:
- **Config**: User configuration
- **Client**: Authenticated SDK client
- **APICtx**: Context with JWT token

### Auth Flow

```
┌─────────────────┐
│ sctl auth login │
└────────┬────────┘
         │
         v
┌──────────────────┐
│ OIDC Device Flow │
│ (Browser Auth)   │
└────────┬─────────┘
         │
         v
┌──────────────────┐
│ JWT Token Saved  │
│ to config.yaml   │
└────────┬─────────┘
         │
         v
┌──────────────────┐
│ All API Calls    │
│ Use JWT Token    │
└──────────────────┘
```

## Development

### Project Structure

```
sctl/
├── cmd/
│   └── main.go              # Entry point with PersistentPreRunE
├── commands/
│   ├── auth.go              # Auth commands
│   ├── config.go            # Config commands
│   ├── workspace.go         # Workspace commands
│   ├── agent*.go            # Agent commands (split)
│   ├── job*.go              # Job commands (split)
│   ├── backup*.go           # Backup commands (split)
│   └── billing*.go          # Billing commands (split)
├── internal/
│   ├── config.go            # Config loading
│   ├── context.go           # CLI context management
│   └── auth.go              # OIDC device flow
└── README.md
```

### Building

```bash
go build -o sctl ./cmd/main.go
```

### Testing

```bash
# Test help
./sctl --help

# Test config
./sctl config init
./sctl config list

# Test auth
./sctl auth status
./sctl auth login

# Test commands (requires authentication)
./sctl workspace list
```

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration
- [fatih/color](https://github.com/fatih/color) - Terminal colors
- [savedhq/sdk-go](https://github.com/savedhq/sdk-go) - API client

## License

MIT
