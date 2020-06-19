#! /usr/bin/env python[3.4.0]
# -*- coding: utf-8 -*-
import os
import sys
import time
from datetime import datetime
import binascii
from coapthon.resources.resource import Resource
#import threading
def Print_value(input_value) :
    hex_value = input_value.split('|')[2]
    value = binascii.unhexlify(hex_value)
    value = str(value)
    #value = value[1:-1]
    print
    print("--------Water Quality-------")
    print("Depth: " + value[:4] + " (m)")
    print("Temp : " + value[4:7] + "(C)")
    print("DO   : " + value[7:11] + "(mg.L)")
    print("BOD  : " + value[11:14] + "(mg.L)")
    print("COD  : " + value[14:17] + "(mg.L)")
    print("SS   : " + value[17:21] + "(mg.L)")
    print("TN   : " + value[21:26] + "(mg.L)")
    print("TP   : " + value[26:31] + "(mg.L)")
    print("TOC  : " + value[31:34] + "(mg.L)")
    print("pH   : " + value[34:])

class Sensor(Resource):
    def __init__(self, name="Sensor", coap_server=None):
        super(Sensor, self).__init__(name, coap_server, visible=True,
                                            observable=True, allow_children=True)
        
        self.payload = "Sensor"
        self.resource_type = "rt1"
        self.content_type = "text/plain"
        self.interface_type = "if1"

    def render_GET(self, request):
        self.payload="Get test" 
        return self

    def render_POST(self, request):
        print("===========CoAP===========")
        print(request)
	Print_value(request.payload.encode('utf-8'))
#        self.producer.send('coap2', key=b'coap', value=bytes(request.payload.encode('utf-8')))
        # self.producer.flush()
    
        return self




        

    
