#!/bin/sh -e
echo DEFAULT_NFSOPTS=${DEFAULT_NFSOPTS} >> /pluginenv
mkdir -p /dockerplugins
if [ -e /run/docker/plugins ]
then
  mount --bind /run/docker/plugins /dockerplugins
fi
mount --rbind /hostcgroup /sys/fs/cgroup
exec /sbin/init
