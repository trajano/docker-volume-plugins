#!/bin/sh -e
docker build -t rootfsimage .
id=$(docker create rootfsimage -h) # id was cd851ce43a403 when the image was created
mkdir -p build/rootfs
docker export "$id" | tar -x -C build/rootfs
docker rm -vf "$id"
docker rmi rootfsimage
cp config.json build
docker plugin create trajano/glusterfs-volume-plugin
