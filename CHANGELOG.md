# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-08-24

### Added
- **Command Alias & Shell Autocorrect**: Added Zsh (`command_not_found_handler`) and Bash (`command_not_found_handle`) shell hook configurations to automatically run `tellme` and show install hints when attempting to run missing registered commands.

### Changed
- **Standardized Caveats Message**: Updated the Homebrew caveats in `.goreleaser.yaml` and installation tips in `README.md` to guide shell profile configurations using clean shell append commands (`>> ~/.zshrc`).

## [1.2.0] - 2026-08-24

### Added
- **Help & Version Flags**: Standard support for help and version flags (both short and long options: `-v`, `--version`, `-version`, `-h`, `--help`, `-help`, `version`, `help`).
- **Dynamic Build Version Injection**: Configured GoReleaser `ldflags` to inject build versions directly into the binary during compilation.

---

## [1.1.0] - 2026-08-24

### Added
- **Environment Profiles (Stack Queries)**: Predefined developer stack profiles (`web`, `mobile`, `backend`, `db`, `devops`, `frontend`) that allow testing all core tools in a stack concurrently (e.g., `tellme web stack`).
- **Typo Tolerance**: Automatic typo detection for the `stack` keyword (e.g., matching `stak`, `stac`, `staack`, etc.) using Levenshtein distance.
- **Dynamic Help Menu**: Informative guide listing available stacks and descriptions when an unrecognized stack name or generic `stack` query is requested.
- **Agent Customizations**: Added `AGENT.md` workspace rules and a development `SKILL.md` runbook under `.agents/` to help automate agent onboarding, testing, and registry contributions.

### Fixed
- **AWS CLI Detection**: Fixed `awscli` lookup failure and added `aws` as an alias mapping to ensure searches for either term successfully detect the system AWS CLI binary.
- **Cache Invalidation**: Cleared stale lookup caches to allow immediate detection of newly installed binaries.

---

## [1.0.0] - 2026-08-14

### Added
- Initial release of `tellme` offline developer tool inspector.
- Parallel tool scanning and system binary location mapping.
- TTY detection and temporal update estimation.
