#!/bin/bash

echo "this is log output" 

FILE="A.csv"

echo "TIME XXXX" > $FILE
echo "#SIGNAL XXXX" >> $FILE

paste <(seq 1000) <(seq 1000 | shuf)>> $FILE
