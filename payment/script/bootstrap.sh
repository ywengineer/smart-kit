#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=payment
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}