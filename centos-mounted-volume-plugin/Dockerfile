FROM centos:7
RUN yum install -q -y go git epel-release yum-utils
RUN go get github.com/trajano/docker-volume-plugins/centos-mounted-volume-plugin && \
  mv $HOME/go/bin/centos-mounted-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -y go git gcc && \
  yum autoremove -q -y && \
  yum clean all && \
  rm -rf /var/cache/yum && \
  find /var/log -type f -delete
