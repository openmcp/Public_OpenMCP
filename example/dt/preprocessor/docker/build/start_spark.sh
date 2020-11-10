#!/bin/bash

./spark-master &
./spark-worker &

spark-submit --master spark://spark-master:7077 --packages org.apache.spark:spark-streaming-kafka-0-8_2.11:2.2.1  --files msgParser.py ./KETI_Preprocessor.py
