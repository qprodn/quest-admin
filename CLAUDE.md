# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Initialize environment (install protoc plugins and tools)
make init

# Generate API files from proto files (pb.go, HTTP/GRPC stubs, OpenAPI)
make api

# Generate internal protobuf files
make config

# Generate Wire dependency injection code
make wire

# Generate all files (api, config, wire)
make all

# Build the project
make build

# Run tests
make test

# Run go generate on all packages
make generate
```

## Architecture Overview

This is a **Kratos-based microservice admin management system** using clean architecture with Google Wire for compile-time dependency injection.

### Layer Structure

```
api/           # Protobuf service definitions
cmd/           # Main application entry point
internal/
  ├── conf/    # Configuration (protobuf-based)
  ├── server/  # HTTP/GRPC server setup, middleware
  ├── service/ # API layer - implements protobuf services
  ├── biz/     # Business layer - use cases and domain logic
  └── data/    # Data layer - repositories, database access
pkg/           # Reusable utilities
```

### Dependency Flow

Wire injects dependencies in this order (outer → inner):
1. **server/** → HTTP/GRPC servers, middleware
2. **data/** → Repositories, Bun DB, Redis, transaction manager
3. **biz/** → Use cases (call repositories)
4. **service/** → Service implementations (call use cases)

Each layer: `exported` defines a `ProviderSet` exposing its dependencies to Wire.

### Key Technologies

- **Kratos v2.9.2** - Microservice framework
- **Wire** - Compile-time dependency injection
- **gRPC + HTTP** - Protobuf with HTTP/JSON transcoding via Google API annotations
- **Bun** - PostgreSQL ORM
- **Redis** - Caching and sessions
- **SA-Token-Go** - Authentication/authorization
- **Zap** - Structured logging

### Adding a New Module

1. Create proto file in `api/proto/{module}/v1/`
2. Run `make api` to generate API stubs
3. Implement domain entities in `internal/biz/{module}/`
4. Implement repository interface in `internal/biz/{module}/` (e.g., `*Repo` interface)
5. Implement repository in `internal/data/{module}/`
6. Implement use case in `internal/biz/{module}/` (e.g., `*Usecase`)
7. Implement service in `internal/service/{module}/`
8. Add providers to respective `ProviderSet` files
9. Run `make wire` to regenerate DI code

### Wire Configuration

The Wire setup in `cmd/quest-admin/wire.go` follows this pattern:

```go
wire.Build(
    server.ProviderSet,
    data.ProviderSet,
    biz.ProviderSet,
    service.ProviderSet,
    newApp
)
```

When adding new dependencies, add them to the appropriate layer's `ProviderSet` and run `make wire`.

### ID Generation

Use `idgen.IDGenerator.NextID(id.{PREFIX})` for distributed unique IDs. Prefixes are defined in `internal/data/idgen/id.go`.

### Transaction Management

Use the transaction manager for cross-repository operations:

```go
err := uc.tm.InTx(ctx, func(ctx context.Context) error {
    // Multiple repository calls here
    return uc.repo1.Create(ctx, ...)
})
```

### Context Utilities

- `ctxs.GetTenantID(ctx)` - Get tenant ID from context
- `ctxs.SetTenantID(ctx, id)` - Set tenant ID in context

### Error Handling

Use `errorx.Err(errkey.Err*)` directly when returning errors:

```go
if existing != nil {
    return errorx.Err(errkey.ErrUserNotFound)
}
```

Error keys are defined in `types/errkey/*.go` with `errorx.Register()`:
```go
var ErrUserNotFound errorx.ErrorKey = "USER_NOT_FOUND"

func init() {
    errorx.Register(ErrUserNotFound, 404, "USER_NOT_FOUND", "user not found")
}
```

### Configuration

Configuration is YAML-based in `configs/` directory. The `env.active` field selects the environment. Config structure is defined in `internal/conf/conf.proto`.
