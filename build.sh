#!/bin/bash
BASENAME=${PWD##*/}
EXPORT="export"
filename=$BASENAME

case $GOOS in
 "windows")     filename+=".exe";;
esac

go build -o $filename

rm -rf $EXPORT/$BASENAME

mkdir -p $EXPORT/$BASENAME
mkdir -p $EXPORT/$BASENAME/configs
cp configs/*.toml $EXPORT/$BASENAME/configs
mv $filename $EXPORT/$BASENAME/