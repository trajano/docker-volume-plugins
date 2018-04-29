Changelog
=========

## 1.3.1

* Used centos:7 as the base for glusterfs client (Fedora mirrors used by gluster/glusterfs-client get very slow)
* Switched all the plugins to "local" scope capability.

## 1.3.0

* NFS volume plugin added
* Fixed issue with volumes being lost on restart
* Dropped the "v" prefix

## v1.2.0

* CentOS Managed volume plugin added
* Fixed security issue with CIFS volume plugin

## v1.1.0

* CIFS volume plugin added
* Refactored code to move all the common parts into a separate package

## v1.0.0

* Initial release, just glusterfs
