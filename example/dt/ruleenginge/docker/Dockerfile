FROM maven:3-jdk-8

ENV HOME=/app

RUN mkdir -p $HOME

WORKDIR $HOME

# 1. add pom.xml only here

ADD . $HOME
#ADD loop_container.py $HOME
# 2. start downloading dependencies

RUN ["/usr/local/bin/mvn-entrypoint.sh", "mvn", "verify", "clean", "--fail-never"]

# 3. add all source code and start compiling
RUN ["mvn", "package", "-DskipTests", "shade:shade"]

EXPOSE 8005

RUN apt-get -o Acquire::Check-Valid-Until=false update && apt-get -y install software-properties-common && apt-get -y install python3.4 && apt-get install -y python3-dev && apt-get install -y python3-setuptools && apt-get install -y python3-pip
RUN pip3 install -r requirements.txt

CMD ["./start.sh"]

