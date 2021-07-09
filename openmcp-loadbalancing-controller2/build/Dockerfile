FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openmcp-loadbalancing-controller2 \
    USER_UID=1001 \
    USER_NAME=openmcp-loadbalancing-controller2

COPY GeoLite2-City.mmdb /root/
COPY dbip-city-lite-2021-07.mmdb /root/

# install operator binary
COPY _output/bin/openmcp-loadbalancing-controller2 ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

#USER ${USER_UID}
