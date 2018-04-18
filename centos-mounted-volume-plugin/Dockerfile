FROM centos/systemd
RUN yum install -q -q -y go git epel-release yum-utils && yum makecache fast
COPY centos-mounted-volume-plugin.service /usr/lib/systemd/system/
COPY init.sh ps.sh /
RUN ln -s /usr/lib/systemd/system/centos-mounted-volume-plugin.service /etc/systemd/system/multi-user.target.wants/centos-mounted-volume-plugin.service && \
  chmod 700 init.sh ps.sh
RUN go get github.com/trajano/docker-volume-plugins/centos-mounted-volume-plugin && \
  mv $HOME/go/bin/centos-mounted-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -q -y go git gcc && \
  yum autoremove -q -q -y && \
  yum clean all && \
  rm -rf /var/cache/yum /etc/mtab && \
  find /var/log -type f -delete
ENV PACKAGES=nfs-utils
ENV MOUNT_TYPE=nfs
CMD [ "/init.sh" ]
# git pull && ./build.sh &&   docker plugin set trajano/centos-mounted-volume-plugin       PACKAGES=nfs-utils       MOUNT_TYPE=nfs       MOUNT_OPTIONS=hard,proto=tcp,nfsvers=4,intr   &&  docker plugin enable trajano/centos-mounted-volume-plugin