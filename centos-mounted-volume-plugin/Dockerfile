FROM centos/systemd AS base

RUN yum install -q -q -y git epel-release yum-utils rsyslog dbus && yum makecache fast && systemctl enable rsyslog.service && \
    curl --silent -L https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz | tar -C /usr/local -zxf - && \
    rm -rf /var/cache/yum /etc/mtab

COPY centos-mounted-volume-plugin/centos-mounted-volume-plugin.service /usr/lib/systemd/system/
COPY centos-mounted-volume-plugin/init.sh /

RUN ln -s /usr/lib/systemd/system/centos-mounted-volume-plugin.service /etc/systemd/system/multi-user.target.wants/centos-mounted-volume-plugin.service && \
    chmod 644 /usr/lib/systemd/system/centos-mounted-volume-plugin.service && \
    chmod 700 /init.sh

FROM base as dev

RUN yum install -q -q -y git && \
    curl --silent -L https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz | tar -C /usr/local -zxf - && \
    rm -rf /var/cache/yum

COPY centos-mounted-volume-plugin/ /root/go/src/github.com/trajano/docker-volume-plugins/centos-mounted-volume-plugin
COPY mounted-volume/ /root/go/src/github.com/trajano/docker-volume-plugins/mounted-volume

RUN cd /root/go/src/github.com/trajano/docker-volume-plugins/centos-mounted-volume-plugin && \
    /usr/local/go/bin/go get

FROM base

COPY --from=dev /root/go/bin/centos-mounted-volume-plugin /

CMD [ "/init.sh" ]
