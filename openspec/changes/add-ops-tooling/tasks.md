# Tasks: Add Operational Tools

- [ ] API Documentation <!-- id: 0 -->
  - [ ] Install `swag` tool and dependencies <!-- id: 1 -->
  - [ ] Add `// @title`, `// @version` comments to `main.go` and `internal/api/router.go` <!-- id: 2 -->
  - [ ] Add annotation comments to Health Check handler as a proof-of-concept <!-- id: 3 -->
  - [ ] Update `Makefile` to include `gen-docs` target running `swag init` <!-- id: 4 -->
  - [ ] Mount `gin-swagger` middleware in `internal/api/router.go` <!-- id: 5 -->
  - [ ] Verify `/swagger/index.html` loads locally <!-- id: 6 -->

- [ ] Configuration Validation <!-- id: 7 -->
  - [ ] Refactor `internal/config` to expose a `Validate()` function separate from `Load()` <!-- id: 8 -->
  - [ ] Implement `validate` logic for Devices (check IDs, Protocols) <!-- id: 9 -->
  - [ ] Implement `validate` logic for Datalinks (check references to Devices/Tasks) <!-- id: 10 -->
  - [ ] Add `validate` subcommand to `cmd/gateway/main.go` using Cobra <!-- id: 11 -->
  - [ ] Wire up the command to call the `Validate()` function <!-- id: 12 -->
  - [ ] Verify `gateway validate` with both valid and invalid configs <!-- id: 13 -->
