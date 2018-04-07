FROM alpine:3.7
RUN apk add -U go git cifs-utils
RUN go get github.com/trajano/docker-volume-plugins/cifs-volume-plugin && \
  mv $HOME/go/bin/cifs-volume-plugin / && \
  rm -rf $HOME/go && \
  yum remove -q -y go git gcc && \
  yum autoremove -q -y
