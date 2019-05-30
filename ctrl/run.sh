#!/bin/sh

nohup ./build/ip-lookup-service__gin &
app_pid=$!

echo "ip-lookup-service__gin running on pid: $app_pid"
echo $app_pid > pid
echo ""
echo "running 'kill \`cat ./pid\`'"
