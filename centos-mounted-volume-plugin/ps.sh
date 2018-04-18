#!/bin/sh
while true
do
ps aux
#mount | grep cgroup
ls -l /var/log /var/run/docker/plugins /run/docker/plugins
sleep 5
done
