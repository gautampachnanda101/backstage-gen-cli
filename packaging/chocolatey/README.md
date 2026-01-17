# Chocolatey Package for backstage-gen-cli

This directory contains the Chocolatey package files for backstage-gen-cli.

## Package Structure

```
chocolatey/
├── backstage-gen-cli.nuspec    # Package specification
├── tools/
│   ├── chocolateyinstall.ps1   # Install script
│   └── chocolateyuninstall.ps1 # Uninstall script
└── README.md
```

## Building the Package

1. Update version and checksum:
   - Update version in `backstage-gen-cli.nuspec`
   - Update checksum in `tools/chocolateyinstall.ps1`

2. Build the package:
   ```powershell
   cd packaging/chocolatey
   choco pack
   ```

3. This creates `backstage-gen-cli.0.1.0-alpha1.nupkg`

## Testing Locally

```powershell
# Install from local package
choco install backstage-gen-cli -s .

# Test
backstage-gen-cli --version
backstage-gen-cli --help

# Uninstall
choco uninstall backstage-gen-cli
```

## Publishing to Chocolatey Community

1. Create account at https://community.chocolatey.org/

2. Get API key from https://community.chocolatey.org/account

3. Push the package:
   ```powershell
   choco apikey --key YOUR_API_KEY --source https://push.chocolatey.org/
   choco push backstage-gen-cli.0.1.0-alpha1.nupkg --source https://push.chocolatey.org/
   ```

4. Wait for moderation (can take a few days)

## Updating the Package

After each release:

1. Update version in `backstage-gen-cli.nuspec`
2. Update URL and checksum in `tools/chocolateyinstall.ps1`
3. Update release notes URL if needed
4. Build and push new package

## Getting Checksums

```powershell
# Download binary and compute checksum
Invoke-WebRequest -Uri "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v0.1.0-alpha.1/backstage-gen-cli-windows-amd64.exe" -OutFile "backstage-gen-cli.exe"
Get-FileHash .\backstage-gen-cli.exe -Algorithm SHA256
```

Or use the .sha256 file from the release:
```powershell
Invoke-WebRequest -Uri "https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v0.1.0-alpha.1/backstage-gen-cli-windows-amd64.exe.sha256" -OutFile "checksum.txt"
Get-Content .\checksum.txt
```

## Automated Updates

Consider using AU (Automatic Updater):
https://github.com/majkinetor/au

```powershell
# update.ps1
import-module au

function global:au_SearchReplace {
    @{
        "tools\chocolateyinstall.ps1" = @{
            "(url64\s*=\s*)('.*')" = "`$1'$($Latest.URL64)'"
            "(checksum64\s*=\s*)('.*')" = "`$1'$($Latest.Checksum64)'"
        }
    }
}

function global:au_GetLatest {
    $releases = Invoke-RestMethod "https://api.github.com/repos/gautampachnanda101/backstage-gen-cli/releases/latest"
    $version = $releases.tag_name -replace '^v', ''
    $url64 = $releases.assets | Where-Object { $_.name -eq "backstage-gen-cli-windows-amd64.exe" } | Select-Object -ExpandProperty browser_download_url

    @{
        Version = $version
        URL64 = $url64
    }
}

update
```
