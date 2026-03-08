---
paths:
  - "sctl/**"
---

# sctl

- CLI tool for interacting with the platform
- Auth: Auth0 user database
- Built with Cobra (CLI) + Viper (config)
- Entry point: `sctl/cmd/main.go`
- Commands in `sctl/commands/` (9 subdirectories)
- Consumes Go SDK (`github.com/savedhq/sdk-go`)
