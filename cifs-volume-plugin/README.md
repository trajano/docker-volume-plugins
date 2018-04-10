CIFS Volume Plugin
======================

This is a managed Docker volume plugin to allow Docker containers to access CIFS shares.  The cifs-utils do not need to be installed on the host and everything is managed within the plugin.  This is intended as a replacement for docker-volume-netshare specifically for CIFS on my organization's swarm.

### Caveats:

- This is a managed plugin only, no legacy support.
- The contents of `/root` is exposed to the plugin as a *read-only* mount so it may access the credential files.
- There is no validation that the credential file itself is secure.
- There is no support for volumes with the `@` symbol as it is used as the escape character.  I may add `@@` as an escape in the future if needed.
- `.netrc` is not used because by it does not support `domain` and it causes extra errors when using `curl` on the system.
- There are many possible options for `mount.cifs` so rather than restricting what can be done, the plugin expects the configuration file to provide all the necessary information (except for credentials).
- If the credential file is present for the share it will append the `credentials=filename` even if the credentials are already part of the options.
- In order to properly support versions use `--alias` when installing the plugin.
- It uses the same format as docker-volume-netshare for the mount points to facilitate migrations.
- **There is no robust error handling.  So garbage in -> garbage out**

## Credentials

Unlike glusterfs, credentails are generally required to access a CIFS share unless it had allowed guest access.  Credentials must be stored in the node under a path in the `/root/` folder.  By default it is `/root/credentials`

To prevent excess quoting, the '@' sign is used as a path separator and will be translated to `/` when trying to process it.  For example the shared volume `foohost/path/subdir` should have a credential file named `foohost@path@subdir`.  It is expected that the file is readable only by root and that is left to the user.

The credential file must have LF line endings and the format is:

    username=value
    password=value
    domain=value 

To protect it within the plugin, the `/root` mount is remounted as a `1m tmpfs` until the credential file is needed (which is on the `Mount` call in which case `/root` is unmounted, the credential file is used then the `tmpfs` is remounted).  It also means that mounting cannot be done in parallel so it will slow down the startup if there are many shares.

### Load order

It is likely that a single share may have multiple subpaths.  Or there's a global default.  For this situation the lookup logic goes as follows given the example of `foohost/path/subdir` and the default credentials path, it will look for the following files in the following order 

1. `/root/credentials/foohost@path@subdir
2. `/root/credentials/foohost@path
3. `/root/credentials/foohost
4. `/root/credentials/default


## Usage

This uses the `driver_opts.cifsopts` to define the list of options to pass to the mount command (a map couldn't be used as some options have no value and will limit future options from being added if I chose to add them.   In addition, the plugin variable `DEFAULT_CIFSOPTS` can be used to set up the default value for `driver_opts.cifsopts` if it is not specified.  For the most part my SMB shares are on Windows and so my `DEFAULT_CIFSOPTS=vers=3.02,mfsymlinks,file_mode=0666,dir_mode=0777`

The `credentials` should not be passed in and will be added automatically if the credentials file is found.  The `volumes.x.name` specifies the host and share path (do not add the `//` it will automatically be added).

Example in docker-compose.yml assuming the alias was set as `cifs`:

    volumes:
      sample:
        driver: cifs
        driver_opts:
          cifsopts: vers=3.02,mfsymlinks,file_mode=0666,dir_mode=0777
        name: "host/share"

The values above correspond to the following mounting command:

    mount -t cifs \
      -o vers=3.02,mfsymlinks,file_mode=0666,dir_mode=0777,credentials=/root/credentials/host@share
      //host/share [generated_mount_point]

## Testing outside the swarm

This is an example of mounting and testing a store outside the swarm.  It is assuming the share is called `noriko/s`.

    docker plugin install trajano/cifs-volume-plugin --grant-all-permissions
    docker plugin enable trajano/cifs-volume-plugin
    docker volume create -d trajano/cifs-volume-plugin --opt cifsopts=vers=3.02,mfsymlinks,file_mode=0666,dir_mode=0777 noriko/s
    docker run -it -v noriko/s:/mnt alpine
