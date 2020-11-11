#!/bin/bash
#./common.sh
spark-submit --master yarn --num-executors 2 --executor-cores 1 --driver-memory 2g --executor-memory 1g --jars spark-streaming-kafka-0-8-assembly_2.11-2.4.1.jar --files msgParser.py KETI_Preprocessor.py

