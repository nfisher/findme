#!/bin/sh -exu

PLAT=`uname -s`
ARCH=`uname -m`

export GOMAXPROCS=2

./consul/${PLAT}/${ARCH}/consul agent -server -bootstrap-expect 1 -data-dir consul/data -ui-dir consul/ui
CONSUL=$!
