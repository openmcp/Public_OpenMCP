# import json
# import os
# import sys
# import requests
# import bs4
# import time
# from kafka import KafkaProducer
# import paho.mqtt.client as mqtt
# import multiprocessing
# from datetime import datetime
# from http.server import BaseHTTPRequestHandler, HTTPServer
# from coapthon.server.coap import CoAP
# from CoapReceiverResources import Sensor

# sys.path.append(os.path.dirname(os.path.abspath(os.path.dirname(__file__))))


# import json

# def searchConfig():
#     return os.environ["DT_CONF_HOME"]+"/configure"
    
# configFilePath = searchConfig()

# open_json = open(configFilePath, 'r')
# datas_result = json.load(open_json)

# IOT_GATEWAY = datas_result['IoTGateway_ip'] 
# HOST_NAME = datas_result['namenode_ip']
# KAFKA_HOST_NAMES = datas_result['kafkaservers_ip']

# KAFKA_PORT = datas_result['kafka_port']
# HTTP_PORT = datas_result['http_port']
# MQTT_PORT = datas_result['mqtt_port']
# COAP_PORT = datas_result['coap_port']

# MQTT_TIMEOUT =datas_result['mqtt_timeout']


# KAFKA_HOST_NAMES_LIST = KAFKA_HOST_NAMES.split(',')

# BOOTSTRAP_SERVERS = ""
# for i in range(len(KAFKA_HOST_NAMES_LIST)):
#         BOOTSTRAP_SERVERS = BOOTSTRAP_SERVERS + KAFKA_HOST_NAMES_LIST[i] + ":" + str(KAFKA_PORT) + ","
# BOOTSTRAP_SERVERS = BOOTSTRAP_SERVERS[:-1]
# print("Kafka Server : " +BOOTSTRAP_SERVERS)


# class MyHandler(BaseHTTPRequestHandler):
    
#     def do_HEAD(s):
#         s.send_response(200)
#         s.send_header("Content-type", "text/html")
#         s.end_headers()

#     def do_GET(s):
#         """Respond to a GET request."""
#         s.send_response(200)
#         s.send_header("Content-type", "text/html")
#         s.end_headers()
#         s.wfile.write(bytes("<html><head><title>HTTP DATA Generator</title></head>\n", "utf-8"))
#         s.wfile.write(bytes("</body></html>\n", "utf-8"))
#         # print("--------HTTP Send--------")
#         # print("SendComplete")
       
#     def do_POST(s):
        
#         producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS)#, api_version=(0,10))
#         content_length = int(s.headers['Content-Length'])
#         post_data = s.rfile.read(content_length)
#         print("===========HTTP===========")
#         print(s.headers)
#         s.send_response(200)
#         s.send_header("Content-type", "text/html")
#         s.end_headers()
#         # s.wfile.write(bytes("<html><body><h1>hi!</h1></body></html>","utf-8"))
#         print("Data : "+str(post_data)[:20]+" .... "+str(post_data)[-20:])
#         producer.send('http2', key=b'http', value=post_data).get()
#         print("HTTP to Kafka Send Complete")
#         # producer.flush()



# def HTTP_Receiver():
#     httpd = HTTPServer((IOT_GATEWAY, HTTP_PORT), MyHandler)
#     print ("MQTT Server Start - " + IOT_GATEWAY)
#     try:
#         httpd.serve_forever()
#     except KeyboardInterrupt:
#         pass
#     httpd.server_close()


# def MQTT_Receiver():

#     producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS)#, api_version=(0,10))

#     client = mqtt.Client()
#     client.connect(IOT_GATEWAY, MQTT_PORT, MQTT_TIMEOUT)

#     def on_connect(client, userdata, flags, rc): 
#           print ("MQTT Server Start - " + IOT_GATEWAY)
#           #print("Connected with result code "+str(rc))
#           client.subscribe("mqtt", qos=0)

#     def on_message(client, userdata, msg):

#           print("===========MQTT===========")
#           print("State : " + str(msg.state))
#           print("Timestamp : " + str(msg.timestamp))
#           print("Dup : " + str(msg.dup))
#           print("Mid : " + str(msg.mid))
#           print("Topic : "+ str(msg.topic))
#           print("QoS : " + str(msg.qos))
#           print("Retain : " + str(msg.retain))
#           print("Payload : "+str(msg.payload)[:20]+" .... "+str(msg.payload)[-20:])
#           #print(msg.payload)
#           print("MQTT to Kafka Send Complete")
#           producer.send('mqtt2', key=b'mqtt', value=msg.payload)
#           # producer.flush()

#     client.on_connect = on_connect
#     client.on_message = on_message

#     client.loop_forever()

    

# class CoAPServer(CoAP):
#     def __init__(self, host, port, producer):
#         CoAP.__init__(self, (host, port))
#         self.add_resource('Sensor/', Sensor(producer))

# def  COAP_Receiver():
#     producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS)#, api_version=(0,10))
#     server = CoAPServer(IOT_GATEWAY, COAP_PORT, producer)
#     try:
#         print ("CoAP Server Start - " + IOT_GATEWAY)
#         server.listen(10)
#     except KeyboardInterrupt:
#         print ("CoAP Server Shutdown")
#         server.close()
#         print ("CoAP Exiting...")

# if __name__ == '__main__':

#     process_mqtt_receiver = multiprocessing.Process(target=MQTT_Receiver)
#     process_mqtt_receiver.daemon = True
#     process_mqtt_receiver.start()

#     process_http_receiver = multiprocessing.Process(target=HTTP_Receiver)
#     process_http_receiver.daemon = True
#     process_http_receiver.start()

#     process_coap_receiver = multiprocessing.Process(target=COAP_Receiver)
#     process_coap_receiver.daemon = True
#     process_coap_receiver.start()

#     process_http_receiver.join()
#     process_mqtt_receiver.join()
#     process_coap_receiver.join()
    




#! /usr/bin/env python[3.4.0]
# -*- coding: utf-8 -*-
import os
import sys
import requests
import bs4
import paho.mqtt.client as mqtt
import multiprocessing
from coapthon.server.coap import CoAP
from coapthon.resources.resource import Resource
import json
import os
from kafka import KafkaProducer

if sys.version_info >= (3, 0):
    from http.server import BaseHTTPRequestHandler, HTTPServer
else:
    from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer



# Config File Load
#def searchConfig():
#    return os.environ["DT_CONF_HOME"]+"/configure"
    
#configFilePath = searchConfig()

#open_json = open(configFilePath, 'r')
#datas_result = json.load(open_json)

#IoTGateway_ip = datas_result['IoTGateway_ip']
#coap_port = datas_result['coap_port']
#mqtt_port = datas_result['mqtt_port']
#http_port = datas_result['http_port']

IoTGateway_connect = os.environ['IOT_SERVICE_CONNECT']
http_port = int(os.environ['HTTP_PORT'])
mqtt_port = int(os.environ['MQTT_PORT'])
coap_port = int(os.environ['COAP_PORT'])


#mqtt_topic = datas_result['mqtt_topic']
#mqtt_timeout = datas_result['mqtt_timeout']

mqtt_topic = os.environ['MQTT_TOPIC']
mqtt_timeout = int(os.environ['MQTT_TIMEOUT'])


#KAFKA_HOST_NAMES = datas_result['kafkaservers_ip']
#KAFKA_HOST_NAMES_LIST = KAFKA_HOST_NAMES.split(',')
#KAFKA_PORT = datas_result['kafka_port']

#BOOTSTRAP_SERVERS = ""
#for i in range(len(KAFKA_HOST_NAMES_LIST)):
#        BOOTSTRAP_SERVERS = BOOTSTRAP_SERVERS + KAFKA_HOST_NAMES_LIST[i] + ":" + str(KAFKA_PORT) + ","
#BOOTSTRAP_SERVERS = BOOTSTRAP_SERVERS[:-1]
#print(BOOTSTRAP_SERVERS)

BOOTSTRAP_SERVERS = os.environ['KAFKA_CONNECT']

IOT_GATEWAY = IoTGateway_connect
#IOT_GATEWAY = "localhost"

HTTP_PORT = http_port
MQTT_PORT = mqtt_port
COAP_PORT = coap_port

MQTT_TIMEOUT = mqtt_timeout
MQTT_TOPIC = mqtt_topic


class MyHandler(BaseHTTPRequestHandler):

    def do_HEAD(s):
        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()

    def do_GET(s):
        """Respond to a GET request."""
        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()
        s.wfile.write(bytes("<html><head><title>HTTP DATA Generator</title></head>\n", "utf-8"))
        s.wfile.write(bytes("</body></html>\n", "utf-8"))

    def do_POST(s):

        producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS)

        content_length = int(s.headers['Content-Length'])
        post_data = s.rfile.read(content_length)

        print("===========HTTP===========")
        print(s.headers)
        #print("Data : "+str(post_data)[:20]+" .... "+str(post_data)[-20:])
        print(post_data.decode())

        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()

        producer.send('http', key=b'http', value=post_data).get()

def HTTP_Receiver():
    httpd = HTTPServer((IOT_GATEWAY, HTTP_PORT), MyHandler)
    try:
        print ("HTTP Server Start - "+IOT_GATEWAY+ " : "+str(HTTP_PORT))
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()


def MQTT_Receiver():

    client = mqtt.Client()
    client.connect(IOT_GATEWAY, MQTT_PORT, MQTT_TIMEOUT)
    producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS)#, api_version=(0,10))

    def on_connect(client, userdata, flags, rc):
          print ("MQTT Server Start - " + IOT_GATEWAY + " : "+str(MQTT_PORT))
          client.subscribe(MQTT_TOPIC, qos=0)

    def on_message(client, userdata, msg):

          print("===========MQTT===========")
          print("State : " + str(msg.state))
          print("Timestamp : " + str(msg.timestamp))
          print("Dup : " + str(msg.dup))
          print("Mid : " + str(msg.mid))
          print("Topic : "+ str(msg.topic))
          print("QoS : " + str(msg.qos))
          print("Retain : " + str(msg.retain))
          #print("Payload : "+str(msg.payload)[:20]+" .... "+str(msg.payload)[-20:])
          print(msg.payload.decode())

          producer.send('mqtt', key=b'mqtt', value=msg.payload)

    client.on_connect = on_connect
    client.on_message = on_message

    client.loop_forever()

class CoAPServer(CoAP):
    def __init__(self, host, port, producer):
        CoAP.__init__(self, (host, port))
        self.add_resource('Sensor/', Sensor(producer))

def COAP_Receiver():
    producer = KafkaProducer(bootstrap_servers=BOOTSTRAP_SERVERS, api_version=(0,10))
    server = CoAPServer(IOT_GATEWAY, COAP_PORT, producer)
    try:
        print ("CoAP Server Start - " + IOT_GATEWAY + " : "+str(COAP_PORT))
        server.listen(10)
    except KeyboardInterrupt:
        print ("CoAP Server Shutdown")
        server.close()
        print ("CoAP Exiting...")


class Sensor(Resource):
    def __init__(self, producer, name="Sensor", coap_server=None):
        super(Sensor, self).__init__(name, coap_server, visible=True,
                                            observable=True, allow_children=True)
        self.producer = producer

    def render_GET(self, request):
        return self

    def render_POST(self, request):
        print("===========CoAP===========")
        req = request.pretty_print()
        self.producer.send('coap', key=b'coap', value=bytes(request.payload.encode('utf-8')))
        print(req)

        return self

if __name__ == '__main__':

    process_mqtt_receiver = multiprocessing.Process(target=MQTT_Receiver)
    process_mqtt_receiver.daemon = True
    process_mqtt_receiver.start()

    process_http_receiver = multiprocessing.Process(target=HTTP_Receiver)
    process_http_receiver.daemon = True
    process_http_receiver.start()

    process_coap_receiver = multiprocessing.Process(target=COAP_Receiver)
    process_coap_receiver.daemon = True
    process_coap_receiver.start()

    process_http_receiver.join()
    process_mqtt_receiver.join()
    process_coap_receiver.join()
