$binary = "terraform-provider-elasticsearch.exe"
$pluginPath = Join-Path $env:APPDATA "terraform.d" "plugins","registry.terraform.io","estaldo","elasticsearch","0.1","windows_amd64"

go build -o "$binary"

if(-not (Test-Path $pluginPath)) {
    New-Item -Type Directory -Path $pluginPath
}

Copy-Item -Path $binary (Join-Path $pluginPath $binary)
