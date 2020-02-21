#! /bin/bash

export VERSION=v.0.1.0
echo version=%VERSION%

pathdir=$(dirname `pwd`)
export GOPATH_BAK=${GOPATH}
export GOPATH=${GOPATH}:${pathdir};

export GOOS=linux
export GOARCH=amd64

echo start install simulate ...
go install -gcflags " -N -l" simulate

echo make ok

export GOPATH=${GOPATH_BAK}
