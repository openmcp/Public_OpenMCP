FROM java:openjdk-8-jdk

#ENV hadoop_ver 2.9.0
ENV spark_ver 2.2.1


# Get Hadoop from US Apache mirror and extract just the native
# libs. (Until we care about running HDFS with these containers, this
# is all we need.)
#RUN mkdir -p /opt && \
#    cd /opt && \
#    curl https://archive.apache.org/dist/hadoop/common/hadoop-${hadoop_ver}/hadoop-${hadoop_ver}.tar.gz | \
#    tar -zx hadoop-${hadoop_ver}/lib/native && \
#    ln -s hadoop-${hadoop_ver} hadoop && \
#    echo Hadoop ${hadoop_ver} native libraries installed in /opt/hadoop/lib/native


ENV SPARK_HOME /opt/spark

RUN mkdir -p /opt && \
    cd /opt && \
    curl https://archive.apache.org/dist/spark/spark-${spark_ver}/spark-${spark_ver}.tgz | \
        tar -zx && \
    ln -s spark-${spark_ver} spark && \
    cd spark && \
    build/mvn -Pyarn -Phadoop-2.9 -Dhadoop.version=2.9.0 -DskipTests clean package && \
    echo Spark ${spark_ver} installed in /opt

RUN cp /opt/spark/conf/spark-env.sh.template /opt/spark/conf/spark-env.sh
RUN cp /opt/spark/conf/spark-defaults.conf.template /opt/spark/conf/spark-defaults.conf

RUN echo "spark.master spark://spark-master:7077" > /opt/spark/conf/spark-defaults.sh
#RUN echo "spark.driver.extraLibraryPath /opt/hadoop/lib/native" > /opt/spark/conf/spark-defaults.sh
RUN echo "spark.app.id KubernetesSpark" > /opt/spark/conf/spark-defaults.sh

# if numpy is installed on a driver it needs to be installed on all
# workers, so install it everywhere
RUN echo "deb [check-valid-until=no] http://archive.debian.org/debian jessie-backports main" > /etc/apt/sources.list.d/jessie-backports.list

# As suggested by a user, for some people this line works instead of the first one. Use whichever works for your case
# RUN echo "deb [check-valid-until=no] http://archive.debian.org/debian jessie main" > /etc/apt/sources.list.d/jessie.list


RUN sed -i '/deb http:\/\/deb.debian.org\/debian jessie-updates main/d' /etc/apt/sources.list

RUN apt-get -o Acquire::Check-Valid-Until=false update && apt-get -y install software-properties-common && apt-get -y install python3.4 && apt-get install -y python3-dev && apt-get install -y python3-setuptools && apt-get install -y python3-pip

RUN apt-get -y install vim
RUN apt-get install -y python3-numpy
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

ENV PATH $PATH:/opt/spark/bin


ADD build /app
WORKDIR /app
#RUN mv spark /opt/
#RUN mv hadoop /opt/

RUN pip3 install -r requirements.txt

#RUN mv spark-defaults.conf /opt/spark/conf/spark-defaults.conf
#RUN mv log4j.properties /opt/spark/conf/log4j.properties

RUN mv /usr/bin/python /usr/bin/python_old
RUN mv /usr/bin/python3 /usr/bin/python

#COPY src /app

ENV PYTHONUNBUFFERED=0

RUN chmod 775 *.sh spark-master spark-worker

