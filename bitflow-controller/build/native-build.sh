#!/bin/bash
home=`dirname $(readlink -f $0)`
root=`readlink -f "$home/../.."`
go build -o "$home/_output/bin/bitflow-controller" $@ "$root/bitflow-controller/cmd/manager"
