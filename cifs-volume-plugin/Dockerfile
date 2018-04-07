FROM oraclelinux:7-slim
RUN yum install -q -y go git cifs-utils
RUN go get github.com/trajano/docker-volume-plugins/cifs-volume-plugin && \
  mv $HOME/go/bin/cifs-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -y go git gcc && \
  yum autoremove -q -y
