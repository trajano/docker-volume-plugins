Docker Managed Volume Plugins
=============================

**This project is forked at https://github.com/marcelo-ochoa/docker-volume-plugins please submit PR or bugs there.**




This project provides managed volume plugins for Docker to connect to [CIFS](https://github.com/trajano/docker-volume-plugins/tree/master/cifs-volume-plugin), [GlusterFS](https://github.com/trajano/docker-volume-plugins/tree/master/glusterfs-volume-plugin) [NFS](https://github.com/trajano/docker-volume-plugins/tree/master/nfs-volume-plugin).

Along with a generic [CentOS Mounted Volume Plugin](https://github.com/trajano/docker-volume-plugins/tree/master/centos-mounted-volume-plugin) that allows for arbitrary packages to be installed and used by mount.

There are two key labels

* `dev` this is an unstable version primarily used for development testing, do not use it on production.
* `latest` this is the latest version that was built which should be ready for use in production systems.

**There is no robust error handling.  So garbage in -> garbage out**
