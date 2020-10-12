FROM centos/systemd as base

RUN yum install -q -q -y epel-release yum-utils nfs-utils rsyslog dbus && yum makecache fast && systemctl enable rsyslog.service && \
    rm -rf /var/cache/yum /etc/mtab && \
    find /var/log -type f -delete

COPY nfs-volume-plugin/nfs-volume-plugin.service /usr/lib/systemd/system/
COPY nfs-volume-plugin/init.sh /

RUN ln -s /usr/lib/systemd/system/nfs-volume-plugin.service /etc/systemd/system/multi-user.target.wants/nfs-volume-plugin.service && \
    chmod 644 /usr/lib/systemd/system/nfs-volume-plugin.service && \
    chmod 700 /init.sh

FROM base as dev

RUN yum install -q -q -y git && \
    curl --silent -L https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz | tar -C /usr/local -zxf - && \
    rm -rf /var/cache/yum

COPY nfs-volume-plugin/ /root/go/src/github.com/trajano/docker-volume-plugins/nfs-volume-plugin
COPY mounted-volume/ /root/go/src/github.com/trajano/docker-volume-plugins/mounted-volume

RUN cd /root/go/src/github.com/trajano/docker-volume-plugins/nfs-volume-plugin && \
    /usr/local/go/bin/go get

FROM base

COPY --from=dev /root/go/bin/nfs-volume-plugin /

CMD [ "/init.sh" ]
