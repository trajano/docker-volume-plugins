This is a managed plugin only, no legacy support.

### Caveats:

- In order to properly support versions use `--alias` when installing the plugin.
- This only supports one glusterfs cluster per instance use `--alias` to define separate instances
- The value of `SERVERS` is initially blank it needs `docker plugin glusterfs set SERVERS=store1,store2` if it is set then it will be used for all servers and low level options will not be allowed.  Primarily this is to control what the deployed stacks can perform.
- There is no robust error handling.  So garbage in -> garbage out

## Operating modes

There are three operating modes listed in order.

### Just the name

This will rely on `SERVERS` being configured and will use the name as the volume mount.  Sub directory mounts are supported as well.  Since `SERVERS` is set it will not allow any other options from being used.

    volumes:
      sample:
        driver: glusterfs
        name: "volume/subdir"

The value of `name` will be used as the `--volfile-id` and `--subdir-mount`.

### Specify the servers

This uses the `driver_opts.servers` to define a list of servers.  This will not work if `SERVERS` is set

    volumes:
      sample:
        driver: glusterfs
        driver_opts:
          servers: store1,store2
        name: "volume/subdir"

The value of `name` will be used as the `--volfile-id` and `--subdir-mount`.  The values above correspond to the following mounting command:

    glusterfs -s store1 -s store2 --volfile-id=volume \
      --subdir-mount=subdir [generated_mount_point]

### Specify the options

This passes the `driver_opts.glusterfsopts` to the `glusterfs` command followed by the generated mount point

    volumes:
      sample:
        driver: glusterfs
        driver_opts:
          glusterfsopts: "--volfile-server=SERVER --volfile-id=abc --subdir-mount=sub"
        name: "whatever"

The value of `name` will not be used for mounting; the value of `driver_opts.glusterfsopts` is expected to have all the volume connection information.

