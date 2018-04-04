FROM gluster/glusterfs-client
#ENV http_proxy='http://ON34C02699866:3128' \
#  https_proxy='http://ON34C02699866:3128'
RUN yum install -q -y go git
COPY . src
WORKDIR /src
RUN go get -d && go build -i -o /glusterfs-volume-plugin && \
  rn -rf /src $GOPATH && \
  yum remove -y go git && \
  yum autoremove -y
WORKDIR /
