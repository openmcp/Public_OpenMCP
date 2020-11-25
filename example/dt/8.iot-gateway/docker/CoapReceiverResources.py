#! /usr/bin/env python[3.4.0]
# -*- coding: utf-8 -*-
import os
import sys
import time
from datetime import datetime

from coapthon.resources.resource import Resource
#import threading


class Sensor(Resource):
    def __init__(self, producer, name="Sensor", coap_server=None):
        super(Sensor, self).__init__(name, coap_server, visible=True,
                                            observable=True, allow_children=True)
        
        self.payload = "Sensor"
        self.resource_type = "rt1"
        self.content_type = "text/plain"
        self.interface_type = "if1"
        self.producer = producer

    def render_GET(self, request):
        self.payload="Get test" 
        return self

    def render_POST(self, request):
        print("===========CoAP===========")
        print(request)
        self.producer.send('coap2', key=b'coap', value=bytes(request.payload.encode('utf-8')))
        print("CoAP to Kafka Send Complete")
        # self.producer.flush()
    
        return self




        

    
