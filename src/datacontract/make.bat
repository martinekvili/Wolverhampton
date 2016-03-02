@echo off

go install

if %errorlevel% equ 0 gofmt -w datacontract.go