FROM gluster/glusterfs-client
ENV http_proxy='http://ON34C02699866:3128' \
  https_proxy='http://ON34C02699866:3128'
RUN yum install -q -y go git
COPY src/ /go/src/
ENV GOPATH=/go
WORKDIR /go
RUN go get github.com/trajano/glusterfs-volume-plugin
RUN go install github.com/trajano/glusterfs-volume-plugin
# -o /glusterfs-volume-plugin
WORKDIR /
#COPY config.json /
# yum autoremove go