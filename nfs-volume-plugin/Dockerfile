FROM centos/systemd
RUN yum install -q -q -y git epel-release yum-utils nfs-utils rsyslog dbus && yum makecache fast && systemctl enable rsyslog.service && \
    curl --silent -L https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz | tar -C /usr/local -zxf -
COPY nfs-volume-plugin.service /usr/lib/systemd/system/
COPY init.sh /
RUN ln -s /usr/lib/systemd/system/nfs-volume-plugin.service /etc/systemd/system/multi-user.target.wants/nfs-volume-plugin.service && \
    chmod 644 /usr/lib/systemd/system/nfs-volume-plugin.service && \
    chmod 700 /init.sh
RUN /usr/local/go/bin/go get github.com/trajano/docker-volume-plugins/nfs-volume-plugin && \
    mv $HOME/go/bin/nfs-volume-plugin / && \
    rm -rf $HOME/go /usr/local/go && \
    yum remove -q -q -y git && \
    yum autoremove -q -q -y && \
    yum clean all && \
    rm -rf /var/cache/yum /etc/mtab && \
    find /var/log -type f -delete
CMD [ "/init.sh" ]
