NFS Mounted Volume Plugin
=========================

This is a managed Docker volume plugin to allow Docker containers to access an NFS mount without installing NFS on the host.

### Caveats:

- This is a managed plugin only, no legacy support.
- This will mount the `/sys/fs/cgroup` folder of the host and expose it to the plugin *as read only*.  `/sys/fs/cgroup` is needed for systemd integration to work.
- In order to properly support versions use `--alias` when installing the plugin.
- **There is no robust error handling.  So garbage in -> garbage out**

## How to use

In order to not make this a free-for-all, only the `device` option is recognized.  Any mount options need to be set up as part of the plugin.  Multiple copies of the plugin can co-exist with different options under different aliases.

The plugin supports the following settings:

* `DEFAULT_NFSOPTS` this corresponds to the default value `-o` parameter of the `mount` command.  It *will* be treated as a single string so it cannot inject the mount points or devices.

When installinng, it is *recommended* that a PLUGINALIAS is specified so that you would know what it is for and can easily control multiple copies of it.  This can be done in an automated fashion as:

    docker plugin install --alias PLUGINALIAS \
      trajano/nfs-volume-plugin \
      --grant-all-permissions --disable
    docker plugin set PLUGINALIAS DEFAULT_NFSOPTS=hard,proto=tcp,nfsvers=4,intr
    docker plugin enable PLUGINALIAS

Example in docker-compose.yml:

    volumes:
      sample:
        driver: PLUGINALIAS
        driver_opts:
          device: "server1:/share_name"
        name: "whatever_name_you_want"

Which yields the following command

    mount -t nfs -o hard,proto=tcp,nfsvers=4,intr server1:/share_name /generated_mount_point

## Testing outside the swarm

This is an example of mounting and testing a store outside the swarm.  It is assuming the share is called `192.168.1.1:/mnt/routerdrive/nfs`.

    docker plugin install trajano/nfs-volume-plugin --grant-all-permissions
    docker plugin enable trajano/nfs-volume-plugin
    docker volume create -d trajano/nfs-volume-plugin --opt device=192.168.1.1:/mnt/routerdrive/nfs --opt nfsopts=hard,proto=tcp,nfsvers=3,intr,nolock nfsmountvolume
    docker run -it -v nfsmountvolume:/mnt alpine
