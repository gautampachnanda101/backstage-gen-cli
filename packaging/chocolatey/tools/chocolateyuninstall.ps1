$ErrorActionPreference = 'Stop'

$packageName = 'backstage-gen-cli'
$toolsDir = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"

$exePath = "$toolsDir\backstage-gen-cli.exe"

if (Test-Path $exePath) {
    Remove-Item $exePath -Force
    Write-Host "backstage-gen-cli has been uninstalled"
} else {
    Write-Host "backstage-gen-cli executable not found at $exePath"
}
