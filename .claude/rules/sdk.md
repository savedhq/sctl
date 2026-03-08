---
paths:
  - "openapi/**"
  - "sdk-*/**"
---

# SDK & CLI

## OpenAPI
- Spec lives at `openapi/`
- SDKs are generated from the spec — do not hand-edit generated files

## Generated SDKs
- `sdk-go/` — Go SDK
- `sdk-python/` — Python SDK
- `sdk-typescript/` — TypeScript SDK

## Consumers
- `dashboard/` — uses the TypeScript SDK
- `sctl/` — CLI tool, uses the Go SDK
