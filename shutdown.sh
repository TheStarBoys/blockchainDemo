#!/bin/bash
pid=`ps -ef | grep blockchainDemo | grep -v grep | awk '{print $2}'`
if [ -z "$pid" ];
then
    echo "[ not find blockchainDemo pid ]"
else
    echo "find result: $pid "
    kill -9 $pid
fi
