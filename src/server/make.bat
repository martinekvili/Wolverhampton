@echo off

go install

if %errorlevel% equ 0 gofmt -w buildresult.go jobhandler.go jobqueue.go queue.go server_main.go servicecontract.go