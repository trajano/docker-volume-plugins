This is a managed plugin only, no legacy support.

### Caveats:

- In order to properly support versions use `--alias` when installing the plugin.
- This only supports one glusterfs cluster per instance use `--alias` to define separate instances
- The value of `SERVERS` is initially blank it needs `docker plugin glusterfs set SERVERS=store1,store2` to be set before using.

    volumes:
      sample:
        driver: glusterfs
        name: "volume/subdir"

The value of `name` will be used as the volume ID.
