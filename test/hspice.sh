#!/bin/bash

echo "this is log output"

[ "$2" == "err" ] && exit 1
true
