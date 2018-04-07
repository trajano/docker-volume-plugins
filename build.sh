#!/bin/sh -e
build() {
    docker build -t rootfsimage $1
    id=$(docker create rootfsimage -h) # id was cd851ce43a403 when the image was created
    rm -rf build/rootfs
    mkdir -p build/rootfs
    docker export "$id" | tar -x -C build/rootfs
    docker rm -vf "$id"
    docker rmi rootfsimage
    cp $1/config.json build
    docker plugin create trajano/$1 build
}
build glusterfs-volume-plugin
build cifs-volume-plugin
