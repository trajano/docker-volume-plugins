CentOS Mounted Volume Plugin
============================

This is a managed Docker volume plugin to allow Docker containers to access any file system that can be mounted.  It is based on CentOS and can utilize any package that CentOS along with EPEL provides.

### Caveats:

- This is a managed plugin only, no legacy support.
- This is *VERY OPEN* so please be careful when using this.  If there is a specialized driver use that one if possible.
- This requires an *INTERNET* connection to download the additional packages.  However `http_proxy` environment variable can be set.
- This slows down your start up as it will download packages on initialization.
- This will mount the `/root` folder of the host and expose it to the plugin *as read only*.  The purpose of which is to allow it to access files that are meant to be read by the mount command for credentials.
- I had to choose one OS distro to start with, if there is a need for another distro just create an issue or PR.
- In order to properly support versions use `--alias` when installing the plugin.
- **There is no robust error handling.  So garbage in -> garbage out**

## How to use

In order to not make this a free-for-all, only the `device` option is recognized.  Any mount options need to be set up as part of the plugin.  Multiple copies of the plugin can co-exist with different options under different aliases.

The plugin supports the following settings:

* `PACKAGES` this is a *comma* separated list of packages that would be added.
* `MOUNT_OPTIONS` this corresponds to the `-o` parameter of the `mount` command.  It *will* be treated as a single string so it cannot inject the mount points or devices.
* `MOUNT_TYPE` the type of the mount, this corresponed to the `-t` parameter  of the `mount` command
* `http_options` (note lower case) this sets the HTTP options as per https://www.centos.org/docs/5/html/yum/sn-yum-proxy-server.html

When installinng, it is *recommended* that a PLUGINALIAS is specified so that you would know what it is for and can easily control multiple copies of it.  This can be done in an automated fashion as:

    docker plugin install --alias PLUGINALIAS \
      trajano/centos-mounted-volume-plugin \
      --grant-all-permissions --disable
    docker plugin set PLUGINALIAS PACKAGES=nfs-utils
    docker plugin set PLUGINALIAS MOUNT_TYPE=nfs
    docker plugin set PLUGINALIAS MOUNT_OPTIONS=hard,proto=tcp,nfsvers=4,intr
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

This is an example of mounting and testing a store outside the swarm.  It is assuming the server is called `store1` and the volume name is `trajano`.

    docker plugin install trajano/glusterfs-volume-plugin --grant-all-permissions
    docker plugin enable trajano/glusterfs-volume-plugin
    docker volume create -d trajano/glusterfs-volume-plugin --opt servers=store1 trajano
    docker run -it -v trajano:/mnt alpine
