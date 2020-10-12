GlusterFS Volume Plugin
=======================

**This project is forked at https://github.com/marcelo-ochoa/docker-volume-plugins please submit PR or bugs there.**

This is a managed Docker volume plugin to allow Docker containers to access GlusterFS volumes.  The GlusterFS client does not need to be installed on the host and everything is managed within the plugin.

### Caveats:

- Requires Docker 18.03-1 at minimum.
- This is a managed plugin only, no legacy support.
- In order to properly support versions use `--alias` when installing the plugin.
- This only supports one glusterfs cluster per instance use `--alias` to define separate instances
- The value of `SERVERS` is initially blank it needs `docker plugin glusterfs set SERVERS=store1,store2` if it is set then it will be used for all servers and low level options will not be allowed.  Primarily this is to control what the deployed stacks can perform.  The values are the DNS or IP addresses of the Gluster servers you are using.
- **There is no robust error handling.  So garbage in -> garbage out**

## Operating modes

There are three operating modes listed in order of preference.  Each are mutually exclusive and will result in an error when performing a `docker volume create` if more than one operating mode is configured.

### Just the name

This is the *recommended* approach for production systems as it will prevent stacks from specifying any random server.  It also prevents the stack configuration file from containing environment specific servers and instead defers that knowledge to the plugin only which is set on the node level.  This relies on `SERVERS` being configured and will use the name as the volume mount set by [`docker plugin set`](https://docs.docker.com/engine/reference/commandline/plugin_set/).  This can be done in an automated fashion as:

    docker plugin install --alias PLUGINALIAS \
      trajano/glusterfs-volume-plugin \
      --grant-all-permissions --disable
    docker plugin set PLUGINALIAS SERVERS=store1,store2
    docker plugin enable PLUGINALIAS

If there is a need to have a different set of servers, a separate plugin alias should be created with a different set of servers.

Example in docker-compose.yml:

    volumes:
      sample:
        driver: glusterfs
        name: "volume/subdir"

The `volumes.x.name` specifies the volume and optionally a subdirectory mount.  The value of `name` will be used as the `--volfile-id` and `--subdir-mount`.  Note that `volumes.x.name` must not start with `/`.

### Specify the servers

This uses the `driver_opts.servers` to define a comma separated list of servers.  The rules for specifying the volume is the same as the previous section.

Example in docker-compose.yml assuming the alias was set as `glusterfs`:

    volumes:
      sample:
        driver: glusterfs
        driver_opts:
          servers: store1,store2
        name: "volume/subdir"

The `volumes.x.name` specifies the volume and optionally a subdirectory mount.  The value of `name` will be used as the `--volfile-id` and `--subdir-mount`.  Note that `volumes.x.name` must not start with `/`.  The values above correspond to the following mounting command:

    glusterfs -s store1 -s store2 --volfile-id=volume \
      --subdir-mount=subdir [generated_mount_point]

### Specify the options

This passes the `driver_opts.glusteropts` to the `glusterfs` command followed by the generated mount point.  This is the most flexible method and gives full range to the options of the glusterfs FUSE client.  Example in docker-compose.yml assuming the alias was set as `glusterfs`:

    volumes:
      sample:
        driver: glusterfs
        driver_opts:
          glusteropts: "--volfile-server=SERVER --volfile-id=abc --subdir-mount=sub"
        name: "whatever"

The value of `name` will not be used for mounting; the value of `driver_opts.glusterfsopts` is expected to have all the volume connection information.

## Testing outside the swarm

This is an example of mounting and testing a store outside the swarm.  It is assuming the server is called `store1` and the volume name is `trajano`.

    docker plugin install trajano/glusterfs-volume-plugin --grant-all-permissions
    docker plugin enable trajano/glusterfs-volume-plugin
    docker volume create -d trajano/glusterfs-volume-plugin --opt servers=store1 trajano
    docker run -it -v trajano:/mnt alpine
