$ErrorActionPreference = "Stop"
$silverc = "C:\Users\Remco\tools\silverc\silverc.exe"
if (-not (Test-Path $silverc)) {
  throw "missing official silverc.exe (kaspanet/silverscript v1-rc1)"
}
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$hits = Select-String -Path "$root\contracts\v1\KasName.sil" -Pattern "readInputState"
if ($hits) {
  throw "readInputState is forbidden on v1-rc1 (silverscript#234)"
}
& $silverc "$root\contracts\v1\KasName.sil" --constructor-args "$root\contracts\v1\ctor-KasName.json" -o "$root\contracts\v1\KasName.json"
Write-Output "compiled KasName.sil"
