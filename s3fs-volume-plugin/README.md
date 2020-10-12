S3Fs Volume Plugin
=======================

**This project is forked at https://github.com/marcelo-ochoa/docker-volume-plugins please submit PR or bugs there.**

This is a managed Docker volume plugin to allow Docker containers to access S3Fs volumes.  The S3Fs client does not need to be installed on the host and everything is managed within the plugin.

### Caveats:

- Requires Docker 18.03-1 at minimum.
- This is a managed plugin only, no legacy support.
- In order to properly support versions use `--alias` when installing the plugin.
- This only supports one S3Fs cluster per instance use `--alias` to define separate instances
- The value of `AWSACCESSKEYID/AWSSECRETACCESSKEY` is initially blank it needs `docker plugin s3fs set AWSACCESSKEYID=key;docker plugin s3fs set AWSSECRETACCESSKEY=secret` if it is set then it will be used for all buckets and low level options will not be allowed.  Primarily this is to control what the deployed stacks can perform.
- **There is no robust error handling.  So garbage in -> garbage out**

## Operating modes

There are three operating modes listed in order of preference.  Each are mutually exclusive and will result in an error when performing a `docker volume create` if more than one operating mode is configured.

### Just the name

This is the *recommended* approach for production systems as it will prevent stacks from specifying any random server.  It also prevents the stack configuration file from containing environment specific key/secrets and instead defers that knowledge to the plugin only which is set on the node level.  This relies on `AWSACCESSKEYID/AWSSECRETACCESSKEY` being configured and will use the name as the volume mount set by [`docker plugin set`](https://docs.docker.com/engine/reference/commandline/plugin_set/).  This can be done in an automated fashion as:

    docker plugin install --alias PLUGINALIAS \
      trajano/s3fs-volume-plugin \
      --grant-all-permissions --disable
    docker plugin set PLUGINALIAS AWSACCESSKEYID=key
    docker plugin set PLUGINALIAS AWSSECRETACCESSKEY=secret
    docker plugin enable PLUGINALIAS

If there is a need to have a different set of key/secrets, a separate plugin alias should be created with a different set of key/secrets.

Example in docker-compose.yml:

    volumes:
      sample:
        driver: s3fs
        name: "bucket/subdir"

The `volumes.x.name` specifies the bucket and optionally a subdirectory mount.  The value of `name` will be used as the `-o bucket=` and `-o servicepath=` in s3fs fuse mount.  Note that `volumes.x.name` must not start with `/`.

### Specify the s3fs driver opts

This uses the `driver_opts.s3fsopts` to define a comma separated list s3fs options.  The rules for specifying the volume is the same as the previous section.

Example in docker-compose.yml assuming the alias was set as `s3fs`:

    volumes:
      sample:
        driver: s3fs
        driver_opts:
          s3fsopts: nomultipart,use_path_request_style
        name: "bucket/subdir"

The `volumes.x.name` specifies the bucket and optionally a subdirectory mount.  The value of `name` will be used as the `-o bucket=` and `-o servicepath=`.  Note that `volumes.x.name` must not start with `/`.  The values above correspond to the following mounting command:

    s3fs -o nomultipart,use_path_request_style,bucket=bucket,servicepath=subdir [generated_mount_point]

### Specify the options

This passes the `driver_opts.s3fsopts` to the `s3fs` command followed by the generated mount point.  This is the most flexible method and gives full range to the options of the S3Fs FUSE client.  Example in docker-compose.yml assuming the alias was set as `s3fs`:

    volumes:
      sample:
        driver: s3fs
        driver_opts:
          s3fsopts: "nomultipart,use_path_request_style,bucket=bucket,servicepath=subdir"
        name: "whatever"

The value of `name` will not be used for mounting; the value of `driver_opts.s3fsopts` is expected to have all the volume connection information.

## Testing outside the swarm

This is an example of mounting and testing a store outside the swarm.  It is assuming the server is called `store1` and the volume name is `trajano`.

    docker plugin install trajano/s3fs-volume-plugin --grant-all-permissions --disable
    docker plugin set AWSACCESSKEYID=key
    docker plugin set AWSSECRETACCESSKEY=secret
    docker plugin set DEFAULT_S3FSOPTS="nomultipart,use_path_request_style"
    docker plugin enable trajano/s3fs-volume-plugin
    docker volume create -d trajano/s3fs-volume-plugin mybucket
    docker run -it -v mybucket:/mnt alpine

## Testing with Oracle Cloud Object storage

Sample usage Oracle Object Storage in S3 compatibilty mode, replace tenant_id and region_id with a proper value:

    docker plugin install trajano/s3fs-volume-plugin --grant-all-permissions --disable
    docker plugin set AWSACCESSKEYID=key
    docker plugin set AWSSECRETACCESSKEY=secret
    docker plugin set DEFAULT_S3FSOPTS="nomultipart,use_path_request_style,url=https://[tenant_id].compat.objectstorage.[region-id].oraclecloud.com/"
    docker plugin enable trajano/s3fs-volume-plugin
    docker volume create -d trajano/s3fs-volume-plugin mybucket
    docker run -it -v mybucket:/mnt alpine
