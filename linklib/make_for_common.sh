#! /bin/bash

mkdir="src/common"
outdir="commonlinklib"
suffix=".so"

function makelinklib(){
    for file in `ls $1` 
    do
        if [ -d $1"/"$file ] 
        then
            makelinklib $1"/"$file
        elif  [[ "${file##*.}" = "go" ]];
        then
            go build -buildmode=c-shared -o ./outdir/${file%.*}$suffix $file
        else
            echo "not find."
        fi
    done
}

makelinklib $mkdir