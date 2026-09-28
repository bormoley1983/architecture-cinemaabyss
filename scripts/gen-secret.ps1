# Renders src/kubernetes/dockerconfigsecret.yaml with the real value from .env.
# Usage:  .\scripts\gen-secret.ps1
# Output: src/kubernetes/.dockerconfigsecret.rendered.yaml (gitignored)
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot

# Load .env (simple KEY=VALUE lines, '#' comments ignored)
$envFile = Join-Path $root '.env'
if (-not (Test-Path $envFile)) { throw ".env not found at $envFile" }
$values = @{}
Get-Content $envFile | ForEach-Object {
    $line = $_.Trim()
    if ($line -and -not $line.StartsWith('#') -and $line.Contains('=')) {
        $k, $v = $line.Split('=', 2)
        $values[$k.Trim()] = $v.Trim()
    }
}

$b64 = $values['DOCKER_CONFIG_JSON_B64']
if (-not $b64 -or $b64 -eq '__DOCKER_CONFIG_JSON_B64__') {
    throw "DOCKER_CONFIG_JSON_B64 is missing or still a placeholder in .env"
}

# Sanity check: must decode to JSON containing auths.ghcr.io.auth
$decoded = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($b64))
if ($decoded -notmatch '"auths".*"ghcr\.io".*"auth"') {
    throw "DOCKER_CONFIG_JSON_B64 does not decode to the expected docker config JSON"
}

$templatePath = Join-Path $root 'src/kubernetes/dockerconfigsecret.yaml'
$outPath = Join-Path $root 'src/kubernetes/.dockerconfigsecret.rendered.yaml'
$content = Get-Content $templatePath -Raw -Encoding utf8
$rendered = $content.Replace('__DOCKER_CONFIG_JSON_B64__', $b64)
Set-Content -Path $outPath -Value $rendered -Encoding utf8 -NoNewline
Write-Host "Rendered secret written to: $outPath"
Write-Host "Apply with: kubectl apply -f src/kubernetes/.dockerconfigsecret.rendered.yaml"
