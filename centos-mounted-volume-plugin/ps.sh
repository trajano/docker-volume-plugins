#!/bin/sh
while true
do
ps aux
mount | grep cgroup
ls -l /var/log
sleep 5
done
