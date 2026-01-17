# Scoop Bucket for backstage-gen-cli

This directory contains the Scoop manifest for backstage-gen-cli.

## Setting Up the Bucket Repository

To distribute via Scoop, create a separate repository named `scoop-backstage-gen-cli`:

1. Create repository: `https://github.com/gautampachnanda101/scoop-backstage-gen-cli`

2. Copy the manifest:
   ```bash
   cp backstage-gen-cli.json /path/to/scoop-backstage-gen-cli/bucket/
   ```

3. Update SHA256 checksum:
   ```bash
   # Get checksum from release assets
   curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v0.1.0-alpha.1/backstage-gen-cli-windows-amd64.exe.sha256
   ```

4. Push to GitHub

## Installing

Users can then install with:

```powershell
scoop bucket add backstage-gen-cli https://github.com/gautampachnanda101/scoop-backstage-gen-cli
scoop install backstage-gen-cli
```

## Updating the Manifest

After each release:

1. Update `version` in the manifest
2. Compute new SHA256:
   ```powershell
   Get-FileHash backstage-gen-cli-windows-amd64.exe -Algorithm SHA256
   ```
3. Update the `hash` field
4. Commit and push

## Autoupdate

The manifest includes `checkver` and `autoupdate` configurations. To automatically update:

```powershell
# In the bucket repository
scoop update backstage-gen-cli
```

## Repository Structure

```
scoop-backstage-gen-cli/
└── bucket/
    └── backstage-gen-cli.json
```

## Testing Locally

```powershell
# Test installation
scoop install ./backstage-gen-cli.json

# Check version
backstage-gen-cli --version

# Uninstall
scoop uninstall backstage-gen-cli
```
