FROM oraclelinux:7-slim
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
RUN yum install -q -y oracle-gluster-release-el7 && \
    yum install -q -y git glusterfs glusterfs-fuse attr && \
    curl --silent -L https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz | tar -C /usr/local -zxf -
RUN /usr/local/go/bin/go get github.com/trajano/docker-volume-plugins/glusterfs-volume-plugin && \
    mv $HOME/go/bin/glusterfs-volume-plugin / && \
    rm -rf $HOME/go /usr/local/go && \
    yum remove -q -y git && \
    yum clean all && \
    rm -rf /var/cache/yum /var/log/anaconda /var/cache/yum /etc/mtab && \
    rm /var/log/lastlog /var/log/tallylog
