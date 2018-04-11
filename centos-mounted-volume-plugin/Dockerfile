FROM centos/systemd
RUN yum install -q -q -y go git epel-release yum-utils && yum makecache fast
COPY centos-mounted-volume-plugin.service /usr/lib/systemd/system/
#COPY centos-mounted-volume-plugin.service /etc/systemd/system/multi-user.target.wants/
#RUN systemctl enable centos-mounted-volume-plugin.service

#COPY systemctl.py /
# /usr/bin/systemctl
#RUN systemctl enable centos-mounted-volume-plugin.service
#RUN chmod 755 /usr/bin/systemctl && \
#RUN chmod 755 /systemctl.py && \
#  ln -s /usr/lib/systemd/system/centos-mounted-volume-plugin.service /etc/systemd/system/multi-user.target.wants/centos-mounted-volume-plugin.service
#COPY centos-mounted-volume-plugin.service /etc/systemd/system/multi-user.target.wants/
RUN go get github.com/trajano/docker-volume-plugins/centos-mounted-volume-plugin && \
  mv $HOME/go/bin/centos-mounted-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -q -y go git gcc && \
  yum autoremove -q -q -y && \
  yum clean all && \
  rm -rf /var/cache/yum && \
  find /var/log -type f -delete
#CMD [ "/systemctl.py", "--verbose" ]
ENV PACKAGES=nfs-utils
ENV MOUNT_TYPE=nfs
#CMD [ "/centos-mounted-volume-plugin" ]
CMD [ "/usr/sbin/init" ]
#CMD [ "/usr/sbin/init", "--verbose" ]
# git pull && ./build.sh
# ls -l /var/lib/docker/plugins/$(docker plugin ls --no-trunc | grep centos | awk '{ print $1 }')/rootfs/usr/lib/systemd/system
# ls -l /var/lib/docker/plugins/$(docker plugin ls --no-trunc | grep centos | awk '{ print $1 }')/rootfs/etc/systemd/system


#FROM centos:7
#ENV container docker
#RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
#systemd-tmpfiles-setup.service ] || rm -f $i; done); \
#rm -f /lib/systemd/system/multi-user.target.wants/*;\
#rm -f /etc/systemd/system/*.wants/*;\
#rm -f /lib/systemd/system/local-fs.target.wants/*; \
#rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
#rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
#rm -f /lib/systemd/system/basic.target.wants/*;\
#rm -f /lib/systemd/system/anaconda.target.wants/*;
#VOLUME [ "/sys/fs/cgroup" ]
#CMD ["/usr/sbin/init"]
