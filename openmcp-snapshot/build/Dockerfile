FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openmcp-snapshot \
    USER_UID=1001 \
    USER_NAME=openmcp-snapshot

# install operator binary
COPY _output/bin/openmcp-snapshot ${OPERATOR}

COPY bin /usr/local/bin
RUN /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}