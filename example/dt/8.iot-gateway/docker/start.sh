#!/bin/bash
#iptables -A INPUT -p tcp -m tcp --dport 8888 -j ACCEPT
#iptables -A INPUT -p tcp -m tcp --dport 5683 -j ACCEPT
#iptables -A INPUT -p tcp -m tcp --dport 1883 -j ACCEPT
mosquitto -d
sleep 5
#python3 loop_container.py
python3 KETI_IoTGateway.py
