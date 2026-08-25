# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2026-08-24

### Added
- **Agent Customizations**: Added `AGENT.md` workspace rules and a development `SKILL.md` runbook under `.agents/` to help automate agent onboarding, testing, and registry contributions.

### Fixed
- **AWS CLI Detection**: Fixed `awscli` lookup failure and added `aws` as an alias mapping to ensure searches for either term successfully detect the system AWS CLI binary.
- **Cache Invalidation**: Cleared stale lookup caches to allow immediate detection of newly installed binaries.


## [0.2.0] - 2026-08-24


### Added
- **Environment Profiles (Stack Queries)**: Predefined developer stack profiles (`web`, `mobile`, `backend`, `db`, `devops`, `frontend`) that allow testing all core tools in a stack concurrently (e.g., `tellme web stack`).
- **Typo Tolerance**: Automatic typo detection for the `stack` keyword (e.g., matching `stak`, `stac`, `staack`, etc.) using Levenshtein distance.
- **Dynamic Help Menu**: Informative guide listing available stacks and descriptions when an unrecognized stack name or generic `stack` query is requested.

---

## [0.1.0] - 2026-08-14

### Added
- Initial release of `tellme` offline developer tool inspector.
- Parallel tool scanning and system binary location mapping.
- TTY detection and temporal update estimation.
