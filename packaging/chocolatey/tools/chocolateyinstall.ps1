$ErrorActionPreference = 'Stop'

$packageName = 'backstage-gen-cli'
$toolsDir = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"

$url64 = 'https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v0.1.0-alpha.1/backstage-gen-cli-windows-amd64.exe'
$checksum64 = 'PLACEHOLDER_WINDOWS_AMD64_SHA256'
$checksumType64 = 'sha256'

$packageArgs = @{
    packageName    = $packageName
    unzipLocation  = $toolsDir
    url64bit       = $url64
    checksum64     = $checksum64
    checksumType64 = $checksumType64
    fileFullPath   = "$toolsDir\backstage-gen-cli.exe"
}

Get-ChocolateyWebFile @packageArgs

Write-Host "backstage-gen-cli has been installed to $toolsDir"
Write-Host "Run 'backstage-gen-cli --help' to get started"
