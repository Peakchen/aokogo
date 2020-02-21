#! /bin/bash
export VERSION=v.0.1.0
echo version=%VERSION%

pathdir=$(dirname `pwd`)
export GOPATH=${pathdir}

export GOOS=linux
export GOARCH=amd64

echo start install GameServer ...
go install -gcflags " -N -l" GameServer

echo make ok

export GOPATH=${GOPATH_BAK}