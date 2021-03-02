FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openmcp-cluster-manager \
    USER_UID=root \
    USER_NAME=openmcp-cluster-manager

# install operator binary
COPY _output/bin/openmcp-cluster-manager ${OPERATOR}

COPY bin /usr/local/bin

RUN  /usr/local/bin/user_setup

RUN chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
