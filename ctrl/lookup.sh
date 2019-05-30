#!/bin/sh

# target process name
process_name="ip-lookup-service__gin"
if [ "$1" != "" ]
then
  process_name="$1"
fi

pid=$(ps -ef | grep "$process_name" | grep -v grep | awk '{ print $2 }')

if [ "$pid" != "" ]
then
  echo "service \"$process_name\" pid: $pid"
else
  echo "NOT Found process named $process_name"
fi
