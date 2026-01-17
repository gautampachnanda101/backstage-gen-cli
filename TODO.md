# TODO

## Package Manager Distribution

- [x] **Homebrew** - Created Formula in `packaging/homebrew/`
  - [x] Write Formula pointing to GitHub release assets
  - [x] Add SHA256 checksums for darwin-amd64 and darwin-arm64
  - [ ] Submit to homebrew-core (after stable release)

- [x] **Scoop** (Windows)
  - [x] Create scoop bucket repository
  - [x] Write JSON manifest with autoupdate support
  - [ ] Test installation on Windows

- [x] **Chocolatey** (Windows)
  - [x] Create `.nuspec` package specification
  - [x] Create install/uninstall scripts
  - [ ] Set up package verification
  - [ ] Submit for community review

- [ ] **APT/DEB** (Debian/Ubuntu)
  - [ ] Create debian package workflow
  - [ ] Set up PPA or direct .deb releases

- [ ] **RPM** (Fedora/RHEL)
  - [ ] Create RPM spec file
  - [ ] Add to release workflow

## MCP Integration

- [x] Research [backstage-mcp](https://github.com/p7ayfu77/backstage-mcp) integration (documented in docs/mcp-integration.md)
- [x] Evaluate [@mexl/backstage-plugin-catalog-backend-module-mcp](https://www.npmjs.com/package/@mexl/backstage-plugin-catalog-backend-module-mcp) (documented)
- [ ] Add MCP client support for pushing catalog entities
- [ ] Implement `backstage-gen-cli push` command for MCP servers

## Documentation

- [x] Improve README.md with badges and better examples
- [x] Add CONTRIBUTING.md (docs/contributing.md)
- [x] Create detailed installation guide for each platform (docs/installation.md)
- [x] Add architecture documentation (docs/architecture.md)
- [x] Document LLM configuration options (docs/configuration.md)
- [x] Add examples for CI/CD integration (docs/ci-cd-examples.md)
- [x] Add MCP integration documentation (docs/mcp-integration.md)
- [x] Add testing documentation (docs/testing.md)
- [x] Add getting started guide (docs/getting-started.md)

## Features

- [x] Add `--quiet` flag to suppress debug output
- [x] Add `--verbose` levels (0-3)
- [ ] Support for monorepo detection (multiple catalog files)
- [ ] Add `backstage-gen-cli validate` for remote Backstage instance validation
- [ ] Template support for custom catalog generation
- [ ] Config file generation wizard improvements

## Testing

- [x] Add unit tests for detector package
- [x] Add unit tests for LLM client
- [x] Add unit tests for config package
- [x] Add unit tests for generator package
- [x] Add unit tests for validator package
- [x] Add unit tests for output package
- [ ] Integration tests with mock LLM responses
- [ ] E2E tests in CI for all supported languages

## CI/CD

- [ ] Add goreleaser for more robust releases
- [ ] Automate changelog generation
- [ ] Add semantic versioning automation
- [ ] Create nightly builds

## Helper Scripts

- [x] Create `scripts/docker-litellm.sh` for LLM container management
- [x] Create `scripts/setup.sh` for installation
- [x] Create `scripts/dev-setup.sh` for developer environment

## Completed

- [x] Hybrid language detection (pattern-first, LLM disambiguation)
- [x] Support for 17+ programming languages
- [x] Cross-platform release workflow (Linux, macOS, Windows)
- [x] v0.1.0-alpha.1 release published
- [x] SHA256 checksums for all binaries
- [x] Improved LLM prompt and JSON extraction
- [x] Verbosity control system (quiet, normal, verbose, debug)
- [x] Comprehensive documentation in docs/ folder
- [x] Package manager templates (Homebrew, Scoop, Chocolatey)
- [x] Unit test coverage for all core packages
- [x] Helper scripts for Docker/LiteLLM setup
