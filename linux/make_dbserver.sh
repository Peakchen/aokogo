#! /bin/bash
export VERSION=v.0.1.0
echo version=%VERSION%

pathdir=$(dirname `pwd`)

export GOPATH=${pathdir}
echo GOPATH
export GOOS="linux"
export GOARCH="amd64"

echo start install DBServer ...
go install -gcflags " -N -l" DBServer

echo make ok

