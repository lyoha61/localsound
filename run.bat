@echo off
:: Open port on Windows
powershell -Command "if (-not (Get-NetFirewallRule -DisplayName 'LocalSound Backend' -ErrorAction SilentlyContinue)) { Start-Process powershell -Verb RunAs -ArgumentList \"New-NetFirewallRule -DisplayName 'LocalSound Backend' -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow\" }"

docker compose up -d