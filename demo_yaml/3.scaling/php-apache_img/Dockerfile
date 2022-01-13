#FROM k8s.gcr.io/hpa-example:latest

#RUN sed -i "s/Listen 80/Listen 8000/" /etc/apache2/ports.conf 

#COPY ports.conf /etc/apache2
#COPY a.yaml /etc/apache2

#RUN /etc/init.d/apache2 restart

#EXPOSE 8080

FROM php:7.4-apache
COPY src/ /var/www/html/
RUN sed -i "s/Listen 80/Listen 8080/" /etc/apache2/ports.conf 
