@echo off
setlocal enabledelayedexpansion
set CGO_ENABLED=1
go build -o ./tmp/main.exe ./backend/server/main.go
