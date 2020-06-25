FROM ubuntu:latest
MAINTAINER kch2 "atyx300@keti.re.kr"
RUN apt-get update -y
RUN apt-get install -y python-pip python-dev build-essential
RUN apt-get install -y mosquitto

ENV DEBCONF_NOWARNINGS yes
RUN apt-get install -y curl
RUN apt-get install -y iputils-ping
RUN apt-get install -y vim
RUN apt-get install -y net-tools
RUN apt-get install -y iptables

COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
RUN chmod 755 start.sh
#ENTRYPOINT ["python"]
#CMD ["KETI_IoTGateway.py"]
CMD ["./start.sh"]
