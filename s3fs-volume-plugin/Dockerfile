FROM oraclelinux:7-slim AS base
ARG TINI_VERSION="v0.18.0"

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini

RUN chmod +x /tini
RUN yum install -q -y oracle-epel-release-el7 && \
    yum install -q -y fuse s3fs-fuse attr && \
    rm -rf /var/cache/yum /etc/mtab

FROM base as dev

RUN yum install -q -y git && \
    rm -rf /var/cache/yum && \
    curl --silent -L https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz | tar -C /usr/local -zxf -

COPY s3fs-volume-plugin/ /root/go/src/github.com/trajano/docker-volume-plugins/s3fs-volume-plugin
COPY mounted-volume/ /root/go/src/github.com/trajano/docker-volume-plugins/mounted-volume

RUN cd /root/go/src/github.com/trajano/docker-volume-plugins/s3fs-volume-plugin && \
    /usr/local/go/bin/go get

FROM base

COPY --from=dev /root/go/bin/s3fs-volume-plugin /
