#! /bin/bash

find ../src -name "*.pb.go" |xargs rm -rf
find -name "*.proto" |awk '{print "protoc --go_out=plugins=grpc:../src " $1}' |sh
find ../src -name "*.pb.go" |xargs sed -i "s/,omitempty//g"
