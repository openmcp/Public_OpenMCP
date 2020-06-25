FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/loadbalancing-controller \
    USER_UID=1001 \
    USER_NAME=loadbalancing-controller

# install operator binary
COPY _output/bin/loadbalancing-controller ${OPERATOR}

COPY bin /usr/local/bin
COPY GeoLite2-City.mmdb /root/
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

#USER ${USER_UID}
