FROM oraclelinux:7-slim as base
ARG TINI_VERSION="v0.18.0"

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

RUN yum install -y oracle-gluster-release-el7 && \
    yum install -y glusterfs glusterfs-fuse attr && \
    yum clean all && \
    rm -rf /var/cache/yum /var/log/anaconda /var/cache/yum /etc/mtab && \
    rm /var/log/lastlog /var/log/tallylog

FROM base as dev

RUN curl --silent -L https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz | tar -C /usr/local -zxf - && \
    yum install -y git

COPY glusterfs-volume-plugin/ /root/go/src/github.com/trajano/docker-volume-plugins/glusterfs-volume-plugin
COPY mounted-volume/ /root/go/src/github.com/trajano/docker-volume-plugins/mounted-volume

RUN cd /root/go/src/github.com/trajano/docker-volume-plugins/glusterfs-volume-plugin && \
    /usr/local/go/bin/go get

FROM base

COPY --from=dev /root/go/bin/glusterfs-volume-plugin /
