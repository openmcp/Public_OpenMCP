FROM ubuntu:16.04
MAINTAINER kch2 "atyx300@keti.re.kr"
RUN apt-get update -y
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN apt-get update
RUN apt-get upgrade -y python3

#RUN apt-get install -y python-pip python-dev build-essential
RUN apt-get install -y mosquitto

ENV DEBCONF_NOWARNINGS yes
RUN apt-get install -y curl
RUN apt-get install -y iputils-ping
RUN apt-get install -y vim
RUN apt-get install -y net-tools
RUN apt-get install -y iptables
RUN apt-get install -y git


COPY . /app
WORKDIR /app
#RUN pip install -r requirements.txt
RUN chmod 755 start.sh


RUN pip3 install -r requirements.txt
RUN git clone https://github.com/Tanganelli/CoAPthon3.git
RUN cd CoAPthon3 && python3 setup.py sdist && pip3 install dist/CoAPthon3-1.0.1.tar.gz -r requirements.txt

ENV PYTHONUNBUFFERED=0

#ENTRYPOINT ["python"]
#CMD ["KETI_IoTGateway.py"]
CMD ["./start.sh"]
