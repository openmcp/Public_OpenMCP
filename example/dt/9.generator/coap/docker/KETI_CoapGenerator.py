#! /usr/bin/env python[3.4.0]
#-*- coding: utf-8 -*-

from coapthon.client.helperclient import HelperClient
from coapthon.messages.request import Request
from coapthon.utils import generate_random_token
from coapthon import defines

import time
import sys
import os
from Generate import *


# 시스템 아규먼트 입니다. (업체명, 크레인명)
COMPANY_NAME = os.environ["COMPANY_NAME"]
CRANE_NAME= os.environ["CRANE_NAME"]

# IoT Gateway가 실행되는 서버의 IP주소 입니다.
IOT_GATEWAY_IP = os.environ["IOT_GATEWAY_IP"]

# Coap 통신을 위한 포트 번호, Path 입니다.
COAP_PORT = int(os.environ["COAP_PORT"])
COAP_PATH="Sensor"

# 데이터 전송 주기로 초단위 입니다.
GEN_PERIOD_SEC = int(os.environ["GEN_PERIOD_SEC"])

# Coap Client 종료를 위한 핸들링 함수입니다.
def handler(signum, f) :
    client.stop()
    client.close()
    sys.exit()

client = None


# 메인 함수로써 반복문을 통해 IoT Gateway의 Coap Server와 통신합니다.
def  main():
        global client

        CraneFullName = COMPANY_NAME + "_" + CRANE_NAME

        client = HelperClient(server=(IOT_GATEWAY_IP, COAP_PORT))
        DataString =""

        while True :
           try :
                DataString = RandomGenerate(CraneFullName)

                if DataString != None:
                    # Post 통신으로 Coap 서버로 데이터를 전송
                    #response = client.post(COAP_PATH, payload = DataString, no_response=True)
                    #printResponse(response)

                    #request = coap_post_no_response(client, COAP_PATH, payload=DataString)
                    request = coap_post_response(client, COAP_PATH, payload=DataString)
                    printRequest(request)


                # sleep을 통해 주기적으로 데이터 전송
                time.sleep(GEN_PERIOD_SEC)

           except Exception as e :
                print(str(e))
                time.sleep(1)
                #client.stop()
                
        client.stop()

def coap_post_response(client, path, payload, callback=None, timeout=None, **kwargs):
    request = client.mk_request(defines.Codes.POST, path)
    request.token = generate_random_token(2)
    request.payload = payload

    request.type = defines.Types["ACK"]

    for k, v in kwargs.items():
        if hasattr(request, k):
            setattr(request, k, v)


    client.send_request(request, callback, timeout)
    return request

def coap_post_no_response(client, path, payload, callback=None, timeout=None, **kwargs):
    request = client.mk_request(defines.Codes.POST, path)
    request.token = generate_random_token(2)
    request.payload = payload

    
    request.add_no_response()
    request.type = defines.Types["NON"]

    for k, v in kwargs.items():
        if hasattr(request, k):
            setattr(request, k, v)


    client.send_request(request, callback, timeout, no_response=True)
    return request

def printRequest(response):
        print("----------------------------------------------")
        print("Protocol : CoAP")
        print(response.pretty_print())
        print("State : Post Complete")
        print("----------------------------------------------")


if __name__ == '__main__':
    main()
