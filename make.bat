@echo off

go install github.com/martinekvili/Wolverhampton/datacontract

go install github.com/martinekvili/Wolverhampton/webserver
xcopy webserver\templates\* ..\..\..\..\bin\templates /s /i /y
xcopy webserver\stat\* ..\..\..\..\bin\stat /s /i /y

go install github.com/martinekvili/Wolverhampton/server