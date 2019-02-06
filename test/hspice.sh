#!/bin/bash

echo "this is log output"

[ "$1" == "err" ] && exit 1
true
