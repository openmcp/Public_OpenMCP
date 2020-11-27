# -*- coding: utf-8 -*-
from influxdb import InfluxDBClient
import sys
import os

RULE_TABLE_NAME = "_Rule"
ERROR_LIST_TABLE_NAME = "_ErrorList"
SENSOR_LIST_TABLE_NAME = "_SensorList"

INFLUXDB = os.environ["INFLUXDB"]
INFLUX_PORT = os.environ["INFLUX_PORT"]
INFLUX_USER = os.environ["INFLUX_USER"]
INFLUX_PW = os.environ["INFLUX_PW"]

RULE_DIR = os.environ["RULE_DIR"]

def get_influx_data():
  host = INFLUXDB
  port = INFLUX_PORT
  user = INFLUX_USER
  password = INFLUX_PW
  dbname = "mydb"

  crane_list = []
  crane_rule_table_list = []
  crane_sensor_table_list = []
  crane_error_table_list = []

  client = InfluxDBClient(host, port, user, password, dbname)

  query = 'show measurements'
  results = client.query(query)
  table_list = list(results.get_points(measurement=None))


  for table in table_list:
    if "_Rule" in table["name"]:
      crane_full_name = table["name"].replace("_Rule", "")
      crane_list.append(crane_full_name)

  for crane_full_name in crane_list :
    query = 'SELECT * FROM ' + crane_full_name + RULE_TABLE_NAME
    results = client.query(query)
    rule_list = list(results.get_points(measurement=crane_full_name + RULE_TABLE_NAME))
    crane_rule_table_list.append(rule_list)

  for crane_full_name in crane_list :
    query = 'SELECT * FROM ' + crane_full_name + ERROR_LIST_TABLE_NAME
    results = client.query(query)
    error_list = list(results.get_points(measurement=crane_full_name + ERROR_LIST_TABLE_NAME))
    crane_error_table_list.append(error_list)

  for crane_full_name in crane_list :
    query = 'SELECT * FROM ' + crane_full_name + SENSOR_LIST_TABLE_NAME
    results = client.query(query)
    sensor_list = list(results.get_points(measurement=crane_full_name + SENSOR_LIST_TABLE_NAME))
    crane_sensor_table_list.append(sensor_list)

  return crane_list, crane_rule_table_list, crane_sensor_table_list, crane_error_table_list

def make_error_rule(file, rule_list, crane_full_name, error_list):

  crane_full_name_text = make_table_seperate(crane_full_name)

  for rule_dict in rule_list:
    condition_text = combine_condition(rule_dict)
    errorName = rule_dict["ErrorName"]
    declare_rule = "rule \"" + crane_full_name+"_"+rule_dict["RuleName"] + "\"\n   dialect \"mvel\"\n"
    rule_when = "   when\n        $sensor:Sensor(" + crane_full_name_text + "&& (" +  condition_text+"))\n"
    rule_then = "   then\n        $sensor.set_ruleErrorCode(\""+ find_error_number(errorName, error_list) + "\");\nend\n\n"
    data = declare_rule + rule_when + rule_then
    file.write(data)


def combine_condition(rule_dict):
  flag = False
  return_text = "("
  for condition_i in range(4):
    if "Condition"+str(condition_i) in rule_dict.keys():
      condition = rule_dict["Condition"+str(condition_i)]
      i = 0
      offset = 0
      while i < len(condition):
        if condition[i] == '&':
          flag = True
          return_text += make_condition(condition[offset:i]) + " && "
          offset = i+2
          i+=1
        i+=1
      if flag == True:
        return_text += make_condition(condition[offset:]) + ') || '
      else:
        return_text += make_condition(condition[offset:]) + ") || "

  return_text = return_text[:len(return_text)-4]
  return return_text


def make_condition(text):
  text = text.replace(" ","")
  return_text = ""
  for i in range(len(text)) :
    if text[i] == '!' or text[i] == '<' or text[i] == '>' or text[i] == '=' :
      return_text = "Double.valueOf($sensor.get_dataMap.get(\"" + text[:i] + "\")) "
      if text[i+1] == '=' :
        return_text += text[i:i+2] + " " + text[i+2:]
      else :
        return_text += text[i:i+1] + " " + text[i+1:]
      break
  return return_text



def make_table_seperate(crane_full_name) :
        return_text = "($sensor.get_dataMap.get(\"CraneFullName\") == " + "\"" + crane_full_name + "\") "
        return return_text

# 에러 이름하고 특정 업체 테이블 스캔값  입력
def find_error_number(error_name, error_list) :
    error_num = ""
    for error_dict in error_list:
      if error_dict["ErrorName"] == error_name:
        return error_dict["ErrorNum"]

def make_sensor_range_rule(file, crane_full_name, sensor_list):
  for sensor_dict in sensor_list:
    sensor_name = sensor_dict["SensorName"]
    rule_name = crane_full_name + "_" + sensor_name

    min_value = str(sensor_dict["NormalRangeMin"])
    max_value = str(sensor_dict["NormalRangeMax"])
    crane_full_name_text = make_table_seperate(crane_full_name)

    declare_rule = "rule \"" + rule_name + "\"\nno-loop true\n   lock-on-active true\n   dialect \"mvel\"\n"
    rule_when_min = "   when\n      sensor:Sensor()\n      $sensor :Sensor(" + crane_full_name_text + "&& ((Double.valueOf($sensor.get_dataMap.get(\"" + sensor_name + "\")) < " + min_value
    rule_when_max = ") || (Double.valueOf($sensor.get_dataMap.get(\"" + sensor_name + "\")) > " + max_value + ")))\n"
    rule_when = rule_when_min + rule_when_max
    rule_then = "   then\n      sensor.set_resultMap(\"" + sensor_name + "\",\"X\");\nend\n\n"
    data = declare_rule + rule_when + rule_then
    file.write(data)

def main():
  crane_list, crane_rule_table_list, crane_sensor_table_list, crane_error_table_list = get_influx_data()

  if len(sys.argv) == 1 :
    for index in range(len(crane_list)) :
      f = open(RULE_DIR + '/' + crane_list[index] + '.drl', 'w')
      data = "package com.rules\nimport com.sample.Sensor\n\ndeclare Sensor\n@role(event)\nend\n\n"
      f.write(data)
      make_error_rule(f, crane_rule_table_list[index], crane_list[index], crane_error_table_list[index])
      make_sensor_range_rule(f, crane_list[index], crane_sensor_table_list[index])
      f.close()

  else :
    index = crane_list.index(sys.argv[1])
    f = open(RULE_DIR + '/' + crane_list[index] + '.drl', 'w')
    data = "package com.rules\nimport com.sample.Sensor\n\ndeclare Sensor\n@role(event)\nend\n\n"
    f.write(data)
    make_error_rule(f, crane_rule_table_list[index], crane_list[index], crane_error_table_list[index])
    make_sensor_range_rule(f, crane_list[index], crane_sensor_table_list[index])
    f.close()

if __name__ == "__main__":
  main()
  print("Create Rule File")

