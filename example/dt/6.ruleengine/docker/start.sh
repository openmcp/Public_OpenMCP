#!/bin/bash

#python /root/DigitalTwin/workspace/src/Mkdrl/KETI_Mkdrl.py
#java -classpath target/RuleEngine-1.0.0-SNAPSHOT.jar com.sample.RuleEngine
#java -classpath target/DigitalTwinCEP-1.0.0-SNAPSHOT.jar com.sample.Main
#python loop_container.py
python3 ./Mkdrl/KETI_Mkdrl.py
./KETI_DigitalTwin_CEP.sh
