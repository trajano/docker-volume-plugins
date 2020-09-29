#!/bin/bash -e
error() {
  printf '\E[31m'; echo "$@"; printf '\E[0m'
}

if [[ $EUID -ne 0 ]]; then
    error "This script should be run using sudo or as the root user"
    exit 1
fi

TAG=$1
build() {
    docker plugin rm -f trajano/$1 || true
    docker rmi -f rootfsimage || true
    docker build -t rootfsimage $1
    id=$(docker create rootfsimage true) # id was cd851ce43a403 when the image was created
    rm -rf build/rootfs
    mkdir -p build/rootfs
    docker export "$id" | tar -x -C build/rootfs
    docker rm -vf "$id"
    cp $1/config.json build
    if [ -z "$TAG" ]
    then
        docker plugin create trajano/$1 build
    else
        docker plugin create trajano/$1:$TAG build
        docker plugin push trajano/$1:$TAG
    fi
}
build glusterfs-volume-plugin
build s3fs-volume-plugin
build cifs-volume-plugin
build nfs-volume-plugin
build centos-mounted-volume-plugin
