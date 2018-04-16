FROM gluster/glusterfs-client
RUN yum install -q -y go git
RUN go get github.com/trajano/docker-volume-plugins/glusterfs-volume-plugin && \
  mv $HOME/go/bin/glusterfs-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -y go git gcc && \
  yum autoremove -q -y && \
  yum clean all && \
  rm -rf /var/cache/yum /var/log/anaconda /var/cache/yum /etc/mtab && \
  rm /var/log/README /var/log/dnf*.log /var/log/lastlog /var/log/tallylog /var/log/hawkey.log
