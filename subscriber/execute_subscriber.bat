@echo off

set GO111MODULE=on
set GOPATH=C:\Users\user\go;C:\Users\user\go\control\pkg\mod\github.com\streadway\amqp@v1.0.0;C:\Users\user\go\adaptive-moms\publisher
set GOROOT=C:\Program Files\Go
set CONFPATH=C:\Users\user\go\adaptive-moms\data

rem Compile publisher
c:
cd C:\Users\user\go\adaptive-moms\subscriber
go run main.go
cd C:\Users\user\go\adaptive-moms\subscriber