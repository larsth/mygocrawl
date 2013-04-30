#!/bin/sh

cd log
go build

cd ..


cd conf
go build
cd ..

cd crawl

go build

cd ..

go build 
