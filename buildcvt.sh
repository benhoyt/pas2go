#!/usr/bin/env bash

cd converted
go build -gcflags="-e" 2>errors.txt
cd ..
