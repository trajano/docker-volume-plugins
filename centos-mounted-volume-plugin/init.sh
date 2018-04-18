#!/bin/sh -e
echo PACKAGES=${PACKAGES} >> /pluginenv
echo MOUNT_OPTIONS=${MOUNT_OPTIONS} >> /pluginenv
echo MOUNT_TYPE=${MOUNT_TYPE} >> /pluginenv
echo http_proxy=${http_proxy} >> /pluginenv
cat /pluginenv
#ls -lR /sys/fs/cgroup
#mount -t cgroup cgroup /sys/fs/cgroup
mount --rbind /hostcgroup /sys/fs/cgroup
mount | grep cgroup
/ps.sh &
exec /sbin/init
