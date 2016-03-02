@echo off

cd datacontract
call make.bat
cd ..

cd webserver
call make.bat
cd ..

cd server
call make.bat
cd ..