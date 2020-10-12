FROM oraclelinux:7-slim as base

RUN yum install -q -y cifs-utils && \
    rm -rf /var/cache/yum /etc/mtab

FROM base as dev

RUN yum install -q -y git tar && \
    curl --silent -L https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz | tar -C /usr/local -zxf - && \
    rm -rf /var/cache/yum

COPY cifs-volume-plugin/ /root/go/src/github.com/trajano/docker-volume-plugins/cifs-volume-plugin
COPY mounted-volume/ /root/go/src/github.com/trajano/docker-volume-plugins/mounted-volume

RUN cd /root/go/src/github.com/trajano/docker-volume-plugins/cifs-volume-plugin && \
    /usr/local/go/bin/go get

FROM base

COPY --from=dev /root/go/bin/cifs-volume-plugin /
