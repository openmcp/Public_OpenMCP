#! /usr/bin/env python[3.4.0]
#-*- coding: utf-8 -*-

import time
import sys
import os
import paho.mqtt.client as mqtt
from Generate import *


# 시스템 아규먼트 입니다. (업체명, 크레인명)
COMPANY_NAME = os.environ["COMPANY_NAME"]
CRANE_NAME= os.environ["CRANE_NAME"]

# IoT Gateway가 실행되는 서버의 IP주소 입니다.
IOT_GATEWAY_IP = os.environ["IOT_GATEWAY_IP"]

# MQTT 통신을 위한 포트 번호 입니다.
MQTT_PORT = int(os.environ["MQTT_PORT"])
MQTT_TOPIC = "mqtt"

MQTT_TIMEOUT = 3
# 데이터 전송 주기로 초단위 입니다.
GEN_PERIOD_SEC = int(os.environ["GEN_PERIOD_SEC"])



# 메인 함수로써 반복문을 통해 IoT Gateway의 Mqtt Server와 통신합니다.
def main():
    CraneFullName = COMPANY_NAME + "_" + CRANE_NAME

    # Mqtt 브로커 서버와 연결
    client = mqtt.Client()
    while True:
        try:    
            client.connect(IOT_GATEWAY_IP, MQTT_PORT, MQTT_TIMEOUT)
            break
        except Exception as e:
            print(e)
            time.sleep(1)

    while True:
        try:
            DataString = RandomGenerate(CraneFullName)

            if DataString != None:
                # Mqtt 브로커 서버로 데이터 전송
                mqttMsgInfo = client.publish(MQTT_TOPIC, DataString,qos=0,retain=False)    
                printResponse(mqttMsgInfo)

            # sleep 함수를 통해 주기적으로 데이터를 전송합니다.
            time.sleep(GEN_PERIOD_SEC)
        except Exception as e:
            print(e)
            time.sleep(1)

def printResponse(mqttMsgInfo):
        print("----------------------------------------------")
        print("Protocol : MQTT")
        print("TOPIC : "+ MQTT_TOPIC)
        print("RC : "+ str(mqttMsgInfo.rc))
        print("MID : " + str(mqttMsgInfo.mid))
        print("State : Publish Complete")
        print("----------------------------------------------")


if __name__ == '__main__':
    main()
