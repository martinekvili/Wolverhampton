@echo off

go install

if %errorlevel% equ 0 gofmt -w callbackcontract.go pagehandlers.go webserver_main.go sseventsource.go sseventbroker.go