@echo off

echo Getting go-bindata
go get -u github.com/jteeuwen/go-bindata/...

echo Making bindata.go
%GOPATH%\bin\go-bindata assets\

echo Getting govend
go get -u github.com/govend/govend

echo Getting vendor packages
%GOPATH%\bin\govend
