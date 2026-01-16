# TODO

## Package Manager Distribution

- [ ] **Homebrew** - Create `homebrew-backstage-gen-cli` tap repository
  - [ ] Write Formula pointing to GitHub release assets
  - [ ] Add SHA256 checksums for darwin-amd64 and darwin-arm64
  - [ ] Submit to homebrew-core (after stable release)

- [ ] **Scoop** (Windows)
  - [ ] Create scoop bucket repository
  - [ ] Write JSON manifest with autoupdate support
  - [ ] Test installation on Windows

- [ ] **Chocolatey** (Windows)
  - [ ] Create `.nuspec` package specification
  - [ ] Set up package verification
  - [ ] Submit for community review

- [ ] **APT/DEB** (Debian/Ubuntu)
  - [ ] Create debian package workflow
  - [ ] Set up PPA or direct .deb releases

- [ ] **RPM** (Fedora/RHEL)
  - [ ] Create RPM spec file
  - [ ] Add to release workflow

## MCP Integration

- [ ] Research [backstage-mcp](https://github.com/p7ayfu77/backstage-mcp) integration
- [ ] Evaluate [@mexl/backstage-plugin-catalog-backend-module-mcp](https://www.npmjs.com/package/@mexl/backstage-plugin-catalog-backend-module-mcp)
- [ ] Add MCP client support for pushing catalog entities
- [ ] Implement `backstage-gen-cli push` command for MCP servers

## Documentation

- [ ] Improve README.md with badges and better examples
- [ ] Add CONTRIBUTING.md
- [ ] Create detailed installation guide for each platform
- [ ] Add architecture documentation
- [ ] Document LLM configuration options
- [ ] Add examples for CI/CD integration (GitHub Actions, GitLab CI)

## Features

- [ ] Add `--quiet` flag to suppress debug output
- [ ] Add `--verbose` levels (0-3)
- [ ] Support for monorepo detection (multiple catalog files)
- [ ] Add `backstage-gen-cli validate` for remote Backstage instance validation
- [ ] Template support for custom catalog generation
- [ ] Config file generation wizard improvements

## Testing

- [ ] Add unit tests for detector package
- [ ] Add unit tests for LLM client
- [ ] Integration tests with mock LLM responses
- [ ] E2E tests in CI for all supported languages

## CI/CD

- [ ] Add goreleaser for more robust releases
- [ ] Automate changelog generation
- [ ] Add semantic versioning automation
- [ ] Create nightly builds

## Completed

- [x] Hybrid language detection (pattern-first, LLM disambiguation)
- [x] Support for 17+ programming languages
- [x] Cross-platform release workflow (Linux, macOS, Windows)
- [x] v0.1.0-alpha.1 release published
- [x] SHA256 checksums for all binaries
- [x] Improved LLM prompt and JSON extraction
