#!/bin/sh -e
echo PACKAGES=${PACKAGES} >> /pluginenv
echo MOUNT_OPTIONS=${MOUNT_OPTIONS} >> /pluginenv
echo MOUNT_TYPE=${MOUNT_TYPE} >> /pluginenv
echo http_proxy=${http_proxy} >> /pluginenv
mkdir -p /dockerplugins
if [ -e /run/docker/plugins ]
then
  mount --bind /run/docker/plugins /dockerplugins
fi
mount --rbind /hostcgroup /sys/fs/cgroup
exec /sbin/init
