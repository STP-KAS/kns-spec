$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$cands = @()
if ($env:SILVERC) { $cands += $env:SILVERC }
$cmd = Get-Command silverc -ErrorAction SilentlyContinue
if ($cmd) { $cands += $cmd.Source }
$cands += @(
  (Join-Path $root "tools\silverc.exe"),
  (Join-Path $root "tools\silverc\silverc.exe"),
  (Join-Path $env:USERPROFILE "tools\silverc-v1\bin\silverc.exe"),
  (Join-Path $env:USERPROFILE "tools\silverc\silverc.exe")
)
$silverc = $cands | Where-Object { $_ -and (Test-Path $_) } | Select-Object -First 1
if (-not $silverc) {
  throw "silverc.exe not found. Official v1.0.0: https://github.com/kaspanet/silverscript/releases/tag/v1.0.0  Windows zip SHA256 3e0d660c15a9e7ac90f3960da24d348b076b1891481bfe758db18accc8a102e1  Set SILVERC=path\to\silverc.exe"
}
$hits = Select-String -Path "$root\contracts\v1\KasName.sil" -Pattern "readInputState" |
  Where-Object { $_.Line -notmatch '^\s*//' }
if ($hits) {
  throw "readInputState is forbidden (silverscript#234 closed unmerged on v1.0.0)"
}
& $silverc "$root\contracts\v1\KasName.sil" --constructor-args "$root\contracts\v1\ctor-KasName.json" -o "$root\contracts\v1\KasName.json"
Write-Output "compiled KasName.sil with $silverc"
