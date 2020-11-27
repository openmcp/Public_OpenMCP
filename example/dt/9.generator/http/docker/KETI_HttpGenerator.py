#! /usr/bin/env python[3.4.0]
#-*- coding: utf-8 -*-

import time
import sys
import os
import requests
from Generate import *


# 시스템 아규먼트 입니다. (업체명, 크레인명)
COMPANY_NAME = os.environ["COMPANY_NAME"]
CRANE_NAME= os.environ["CRANE_NAME"]

# IoT Gateway가 실행되는 서버의 IP주소 입니다.
IOT_GATEWAY_IP = os.environ["IOT_GATEWAY_IP"]

# Http 통신을 위한 포트 번호, Path 입니다.
HTTP_PORT = int(os.environ["HTTP_PORT"])

# 데이터 전송 주기로 초단위 입니다.
GEN_PERIOD_SEC = int(os.environ["GEN_PERIOD_SEC"])



# 메인 함수로써 반복문을 통해 IoT Gateway의 Http Server와 통신합니다.
def main():
    CraneFullName = COMPANY_NAME + "_" + CRANE_NAME

    while(True):
        try:
            DataString = RandomGenerate(CraneFullName)

            if DataString != None:
                url = "http://" + IOT_GATEWAY_IP + ":" + str(HTTP_PORT)
                # IoT Gateway로 전송
                r = requests.post(url, data=DataString)
                printResponse(r)

            # sleep 함수를 통해 주기적으로 데이터를 전송합니다.
            time.sleep(GEN_PERIOD_SEC)

        except Exception as e:
            print(e)
            time.sleep(1)

def printResponse(r):
        print("----------------------------------------------")
        print("Protocol : HTTP")
        print("Server : " + r.headers['Server'])
        print("Date : " + r.headers['Date'])
        print("Content-Type : " + r.headers['Content-type'])
        print("Status Code : "+ str(r.status_code))
        print("State : Post Complete")
        print("----------------------------------------------")

if __name__ == '__main__':
    main()
