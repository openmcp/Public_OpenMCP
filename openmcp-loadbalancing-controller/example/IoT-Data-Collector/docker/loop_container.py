import json
import os
import sys
import requests
import bs4
from kafka import KafkaProducer
import paho.mqtt.client as mqtt
import multiprocessing
from datetime import datetime
import binascii
import time

if sys.version_info >= (3, 0):
    from http.server import BaseHTTPRequestHandler, HTTPServer
else:
    from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer

from coapthon.server.coap import CoAP
from CoapReceiverResources import Sensor

IOT_SERVICE_CONNECT = os.environ['IOT_SERVICE_CONNECT']

MQTT_TIMEOUT = int(os.environ['MQTT_TIMEOUT'])
KAFKA_CONNECT = os.environ['SERVER_kafka_connect']

HTTP_PORT = 8888
MQTT_PORT = 1883
COAP_PORT = 5683
KAFKA_PORT = 9092

def Print_value(input_value) :
	hex_value = input_value.split('|')[2]
	value = binascii.unhexlify(hex_value)
	value = str(value)

	print
        print("--------Water Quality-------")
        print("Depth: " + value[:4] + "(m)")
        print("Temp : " + value[4:7] + " (C)")
        print("DO   : " + value[7:11] + " (mg.L)")
        print("BOD  : " + value[11:14] + " (mg.L)")
        print("COD  : " + value[14:17] + " (mg.L)")
        print("SS   : " + value[17:21] + " (mg.L)")
        print("TN   : " + value[21:26] + " (mg.L)")
        print("TP   : " + value[26:31] + " (mg.L)")
        print("TOC  : " + value[31:34] + " (mg.L)")
        print("pH   : " + value[34:])
	print

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
        # print("--------HTTP Send--------")
        # print("SendComplete")
       
    def do_POST(s):
        
        #producer = KafkaProducer(bootstrap_servers=KAFKA_CONNECT)#, api_version=(0,10))
        content_length = int(s.headers['Content-Length'])
        post_data = s.rfile.read(content_length)
        print("===========HTTP===========")
        print(s.headers)
        s.send_response(200)
        s.send_header("Content-type", "text/html")
        s.end_headers()
        # s.wfile.write(bytes("<html><body><h1>hi!</h1></body></html>","utf-8"))
        #print("Data : "+str(post_data)[:20]+" .... "+str(post_data)[-20:])
        print("Data : " +  str(post_data))
	Print_value(str(post_data))
        #producer.send('http2', key=b'http', value=post_data).get()
        #print("HTTP to Kafka Send Complete")
        # producer.flush()


def HTTP_Receiver():
    httpd = HTTPServer((IOT_SERVICE_CONNECT, HTTP_PORT), MyHandler)
    print("HTTP Server Start")
    print("----------------------")
    print("Port Num : 8888")
    print("Node Port Num : 30802")
    print
    try:
        httpd.serve_forever()

    except KeyboardInterrupt:
        pass
    httpd.server_close()



def MQTT_Receiver():

    #producer = KafkaProducer(bootstrap_servers=KAFKA_CONNECT)#, api_version=(0,10))

    client = mqtt.Client()
    client.connect(IOT_SERVICE_CONNECT, MQTT_PORT, MQTT_TIMEOUT)

    def on_connect(client, userdata, flags, rc): 
          print ("MQTT Server Start")
          print ("----------------------")
          print ("Port Num : 1883")
          print ("Node Port Num : 30487")
          print
          #print("Connected with result code "+str(rc))
          client.subscribe("mqtt", qos=0)

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
          print("Payload : " + str(msg.payload))
          Print_value(str(msg.payload))
          #print(msg.payload)
          #print("MQTT to Kafka Send Complete")
     #     producer.send('mqtt2', key=b'mqtt', value=msg.payload)
          # producer.flush()

    client.on_connect = on_connect
    client.on_message = on_message

    client.loop_forever()

    

class CoAPServer(CoAP):
    def __init__(self, host, port): #producer):
        CoAP.__init__(self, (host, port))
        self.add_resource('Sensor/', Sensor())

def  COAP_Receiver():
    #producer = KafkaProducer(bootstrap_servers=KAFKA_CONNECT)#, api_version=(0,10))
    server = CoAPServer(IOT_SERVICE_CONNECT, COAP_PORT)
    try:
        print ("CoAP Server Start")
        print ("----------------------")
        print ("Port Num : 5683")
        print ("Node Port Num : 30402")
        print

        server.listen(10)
    except KeyboardInterrupt:
        print ("CoAP Server Shutdown")
        server.close()
        print ("CoAP Exiting...")
        
if __name__ == '__main__':

   while(True) :
	print("adsf")
	time.sleep(1) 
