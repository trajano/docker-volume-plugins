FROM oraclelinux:7-slim
RUN yum install -q -y git cifs-utils tar && \
  curl --silent -L https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz | tar -C /usr/local -zxf -
RUN /usr/local/go/bin/go get github.com/trajano/docker-volume-plugins/cifs-volume-plugin && \
  mv $HOME/go/bin/cifs-volume-plugin / && \
  rm -rf $HOME/go /usr/local/go && \
  yumdb set reason dep git tar && \
  yum autoremove -y && \
  yum clean all && \
  rm -rf /var/cache/yum /etc/mtab && \
  find /var/log -type f -delete
