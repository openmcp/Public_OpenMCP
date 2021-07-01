FROM ubuntu:16.04

RUN  apt-get -y update
RUN apt-get install -y vim
RUN apt-get install -y sysbench

RUN cp /usr/bin/sysbench /usr/local/bin/sysbench-cpu
RUN cp /usr/bin/sysbench /usr/local/bin/sysbench-memory
RUN cp /usr/bin/sysbench /usr/local/bin/sysbench-disk
RUN cp /usr/bin/sysbench /usr/local/bin/sysbench-network

ENV OPERATOR=/usr/local/bin/sysbench-test \
    USER_UID=1001 \
    USER_NAME=sysbench-test

# install operator binary
COPY _output/bin/sysbench-test ${OPERATOR}

COPY bin /usr/local/bin
#RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

