#!/bin/sh

echo "going to terminate ip-lookup service:"

read -n 1 -p "
Are you SURE? 
- Y (default)
- N (cancel)

" to_terminate

if [ $to_terminate == "Y" ]; then
  kill `cat ./pid`

elif [ $to_terminate == "y" ]; then
  kill `cat ./pid`

elif [ $to_terminate == "" ]; then
  kill `cat ./pid`

fi
