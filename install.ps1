#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

$Owner = '{{ .Owner }}'
$Repo = '{{ .Repo }}'
$Version = '{{ .Version }}'
$Binary = '{{ .Binary }}'

$BinDir = "$Home\bin"

$Target = "${Binary}_windows_amd64.tar.gz"
$TargetZip = "$BinDir\${$Target}"
$BinaryExe = "$BinDir\${Binary}.exe"

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$TargetUri = "https://github.com/${owner}/${repo}/releases/download/${Version}/${$Target}"

# create bin dir if not exist
if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

Invoke-WebRequest $TargetUri -OutFile $TargetZip -UseBasicParsing

if (Get-Command Expand-Archive -ErrorAction SilentlyContinue) {
  Expand-Archive $TargetZip -Destination $BinDir -Force
} else {
  if (Test-Path $BinaryExe) {
    Remove-Item $BinaryExe
  }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [IO.Compression.ZipFile]::ExtractToDirectory($TargetZip, $BinDir)
}

Remove-Item $TargetZip

# add $HOME/bin to the $PATH
$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Write-Output "${Binary} was installed successfully to $BinaryExe"
Write-Output "Run '${Binary} --help' to get started"