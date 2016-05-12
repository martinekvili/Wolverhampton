@echo off

go install

if %errorlevel% equ 0 gofmt -w jobhandler.go jobqueue.go queue.go server_main.go servicecontract.go