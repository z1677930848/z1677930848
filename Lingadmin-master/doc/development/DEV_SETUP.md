# Developer Setup — Lingadmin / EdgeAdmin

This file explains how to prepare a local development environment and common workflows.

Prerequisites
- Go 1.22 installed
- Git
- For builds and packaging: a Unix-like shell (WSL on Windows, macOS or Linux). The provided `build/build.sh` is a bash script.

Recommended layout
- Clone this repository and its sibling repos as siblings (same parent dir):

```bash
# from the parent directory
git clone https://github.com/TeaOSLab/EdgeCommon.git
git clone https://github.com/TeaOSLab/EdgeAPI.git
git clone https://github.com/TeaOSLab/EdgeAdmin.git
```

Why sibling repos?
- `go.mod` uses `replace github.com/TeaOSLab/EdgeCommon => ../EdgeCommon`.
- `build/build.sh` expects `EdgeAPI` under `../EdgeAPI` for packaging.

Quick start (development)
- Run the admin server (PowerShell example):

```powershell
# in repository root
go run .\cmd\lingcdnadmin\main.go
# switch environment
go run .\cmd\lingcdnadmin\main.go dev
```

Naming & consistency
- Historical names: the project and scripts sometimes use multiple names (for historical reasons):
	- `EdgeAdmin` / `edge-admin` — upstream project name.
	- `ling-admin` — service/process name used in scripts and constants (see `internal/const/const.go`).
	- `LingCDN` — product name shown to users/UI.

	When making code changes, prefer using `internal/const/const.go` constants (`ProductName`, `ProcessName`) for programmatic references. Keep user-facing strings (UI) separate from process/service identifiers.

Consistency tips:
- Use `ProcessName` when registering services, creating sockets, or referring to the running process.
- Use `ProductName`/`ProductNameZH` for UI text and product branding.
- Document any new global constants in `internal/const/const.go` with comments explaining purpose and compatibility notes.

Building a release (use WSL or Linux/macOS)

```bash
# from repo root inside WSL or Linux
./build/build.sh linux amd64 community
```

Notes and environment variables
- `EDGEAPI_PATH` — optional environment variable to point to the EdgeAPI repository location. Example:

```bash
EDGEAPI_PATH=/home/you/projects/EdgeAPI ./build/build.sh linux amd64 community
```

- `configs/server.template.yaml` is used to generate `configs/server.yaml` on first run if missing.

Secret handling
- The runtime secret used for CSRF and internal token generation is generated at startup and stored in memory (`configs.Secret`). The code intentionally does not persist this secret to source files.
- Do NOT commit production secrets to the repository. For automation, you can set a stable secret in a secure location and inject it into the process environment or into your configuration management system.

CI notes
- The repository has a lightweight GitHub Actions workflow to run `gofmt` and `go test`. The workflow clones `EdgeCommon` and `EdgeAPI` into the parent directory so `replace` directives work in CI.

Troubleshooting
- If `build/build.sh` fails complaining about missing tools (`zip`, `unzip`, `go`, `find`, `sed`), install them or use WSL.
- If you see runtime errors related to missing sibling repos, ensure `../EdgeCommon` and `../EdgeAPI` exist or set `EDGEAPI_PATH`.

If you want, I can add a PowerShell helper script for Windows-native builds, or a Makefile to unify commands across platforms.
