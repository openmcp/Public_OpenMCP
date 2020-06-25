FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openmcp-ingress-controller \
    USER_UID=1001 \
    USER_NAME=openmcp-ingress-controller

# install operator binary
COPY _output/bin/openmcp-ingress-controller ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
