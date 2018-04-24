FROM centos/systemd
RUN yum install -q -q -y go git epel-release yum-utils nfs-utils rsyslog dbus && yum makecache fast && systemctl enable rsyslog.service
COPY nfs-volume-plugin.service /usr/lib/systemd/system/
COPY init.sh /
RUN ln -s /usr/lib/systemd/system/nfs-volume-plugin.service /etc/systemd/system/multi-user.target.wants/nfs-volume-plugin.service && \
  chmod 644 /usr/lib/systemd/system/nfs-volume-plugin.service && \
  chmod 700 /init.sh
RUN go get github.com/trajano/docker-volume-plugins/nfs-volume-plugin && \
  mv $HOME/go/bin/nfs-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -q -y go git gcc && \
  yum autoremove -q -q -y && \
  yum clean all && \
  rm -rf /var/cache/yum /etc/mtab && \
  find /var/log -type f -delete
CMD [ "/init.sh" ]
