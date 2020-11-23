FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/openmcp-globalcache \
    USER_UID=1001 \
    USER_NAME=openmcp-globalcache

# install operator binary
COPY _output/bin/openmcp-globalcache ${OPERATOR}

COPY bin /usr/local/bin
RUN /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]
USER ${USER_UID}