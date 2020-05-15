FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

RUN microdnf update && microdnf install procps

ENV OPERATOR=/usr/local/bin/openmcp-scheduler \
    USER_UID=1001 \
    USER_NAME=openmcp-scheduler

# install operator binary
COPY _output/bin/openmcp-scheduler ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
