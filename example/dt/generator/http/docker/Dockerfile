FROM ubuntu:16.04
MAINTAINER kch2 "atyx300@keti.re.kr"

RUN apt-get update -y
#RUN apt-get install -y python-pip python-dev build-essential
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN apt-get update 
RUN apt-get upgrade -y python3


ENV DEBCONF_NOWARNINGS yes
RUN apt-get install -y curl
RUN apt-get install -y iputils-ping
RUN apt-get install -y vim
RUN apt-get install -y net-tools
RUN apt-get install -y iptables
RUN apt-get install -y git

COPY . /app
WORKDIR /app

RUN pip3 install -r requirements.txt
RUN chmod 755 start_http.sh

ENV PYTHONUNBUFFERED=0

CMD ["/app/start_http.sh"]

