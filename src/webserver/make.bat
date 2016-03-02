@echo off

go install

if %errorlevel% equ 0 gofmt -w webserver_main.go sseventsource.go sseventbroker.go