from pyspark import SparkContext, SparkConf, SparkFiles
from pyspark.streaming import StreamingContext
from pyspark.streaming.kafka import KafkaUtils
from msgParser import *
import msgParser

import json
import happybase
import sys
import os
import pika
import time

from influxdb import InfluxDBClient



SPARK_PERIOD_SEC = int(os.environ["SPARK_PERIOD_SEC"]) # 5
ZKQUORUM = os.environ['ZKQUORUM'] #"zk-cs.datacenter.svc.cluster.local:2181"

QUEUE_NAME = os.environ["QUEUE_NAME"]
QUEUE_TOPIC = os.environ["QUEUE_TOPIC"]

RABITMQ = os.environ["RABITMQ"]

sensorList_table = "_SensorList"
data_table = "_Data"
version_table = "_Version"


INFLUXDB = os.environ["INFLUXDB_SERVICE"]


def saveInfluxDBData(data_dict) :
  client = InfluxDBClient(INFLUXDB, 8086)
  try:
      client.create_database('mydb')
  except Exception as e:
      print(e)
      pass

  try:
      client.switch_database('mydb')
  except Exception as e:
      print(e)
      pass


  
  data_dict = json.loads(data_dict)
  CraneFullName = data_dict['CraneFullName']
  
  input_value = []

  for key , value in data_dict.items() :
      dict_put = { 'measurement' : CraneFullName + data_table,
                   'tags' : { 'SensorName' : key },
                   'fields' : { 'Value' : value }
                 }
      input_value.append(dict_put)

  client.write_points(input_value)

  

def streaming_set() :
  sc = SparkContext(appName="PythonStreamingPreprocessing")
  ssc = StreamingContext(sc, SPARK_PERIOD_SEC) # 1 second window
  zkQuorum = ZKQUORUM

  #Dict of (topic_name -> numPartitions) to consume. Each partition is consumed in its own thread.
  topics = {'http': 1, 'mqtt' : 1, 'coap' : 1}
  stream = KafkaUtils.createStream(ssc, zkQuorum, "raw-event-streaming-consumer", topics, {"auto.offset.reset": "largest"})
  return ssc, stream



def msg_parse(data) :
  data_list = data.split("|")

  CraneFullName = data_list[0]
  Time = data_list[1]
  OBD_DataList = data_list[2:-1]
  Crane_Data = data_list[-1]

  getDataFromInfluxDB(INFLUXDB, CraneFullName)

  OBD_II_data_dict = {}
  sensor_data_dict = {}

  if Crane_Data == None:
    pass

  else:
    sensor_data_encoded_hex = Crane_Data
    sensor_data_dict = sensor_data_parser(sensor_data_encoded_hex)

    if sensor_data_dict == 'ignore' :
      error_dict = {}
      error_dict["FormatError"] = "Crane"
      error_dict["CraneFullName"] = CraneFullName
      jsonString = makeJson(error_dict)
      return jsonString

  data_dict = make_dict(Time, CraneFullName, OBD_II_data_dict , sensor_data_dict)


  jsonString = makeJson(data_dict)
 
  return jsonString

credentials = pika.PlainCredentials('dtuser01', 'dtuser01')
params = pika.ConnectionParameters(RABITMQ, 5672, '/', credentials)
#queue_connection = pika.BlockingConnection(params)

#channel = queue_connection.channel()
#channel.queue_declare(queue=QUEUE_NAME)


def transfer_list(json_list) :
  global channel
  global queue_connection

  for json_data in json_list:
    check_dict = json.loads(json_data)
    if not "FormatError" in check_dict:
        try:
          channel.basic_publish(exchange='',routing_key=QUEUE_TOPIC, body=json_data)
        except Exception as e:
          print("AMQP publish Exception.." +str(e))
          print("recreate connection...")
          print("JsonData : "+str(json_data))
          queue_connection = pika.BlockingConnection(params)

          channel = queue_connection.channel()
          channel.queue_declare(queue=QUEUE_NAME)
          channel.basic_publish(exchange='',routing_key=QUEUE_TOPIC, body=json_data)

        saveInfluxDBData(json_data)
    else:
       print("ERROR")



def parse_function(rdd) :

  json_format_rdd = rdd.map(lambda data : msg_parse(data))

  json_format_list = json_format_rdd.take(int(json_format_rdd.count()))
  if len(json_format_list) != 0:
    transfer_list(json_format_list)




def main() :
  ssc, stream = streaming_set()

  raw_data = stream.map(lambda value: value[1])
  raw_data.foreachRDD(lambda rdd : parse_function(rdd) )

  ssc.start()
  ssc.awaitTermination()


if __name__ == "__main__":
  main()
  # test_msg = "SHINHAN_Crane_1|2020-11-06 11:09:20|H 135.2000.0-07.6000.0364.9000.001030101999.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9 99.9000.000000000.00000000000000000010000000000"
  # test_msg = "SHINHAN_Crane_1|3|2019-12-19 10:38:14|+287.724852+20.527328+22.54130544.3022.4277.4003266909840113017505822271-07.89-2.28+1.61+0.00-07.27+0.4001112209N"
  # jsonString = msg_parse(test_msg)
  # saveInfluxDBData(jsonString)
 


