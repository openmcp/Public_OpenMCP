# -*- coding: utf-8 -*-
import binascii
import json
import happybase
import os
import datetime
from influxdb import InfluxDBClient
import sys

Description = [0 for i in range(0, 100)]

Description[0] = "PIDs supported [01 - 20]"
Description[1] = "Monitor status since DTCs cleared. (Includes malfunction indicator lamp (MIL) status and number of DTCs.)"
Description[2] = "Freeze DTC"
Description[3] = "Fuel system status"
Description[4] = "Calculated_engine_load"
Description[5] = "Engine_coolant_temperature"
Description[6] = "Short term fuel trim—Bank 1"
Description[7] = "Long term fuel trim—Bank 1"
Description[8] = "Short term fuel trim—Bank 2"
Description[9] = "Long term fuel trim—Bank 2"
Description[10] = "Fuel pressure (gauge pressure)"
Description[11] = "Intake manifold absolute pressure"
Description[12] = "Engine RPM"
Description[13] = "Vehicle speed"
Description[14] = "Timing advance"
Description[15] = "Intake air temperature"
Description[16] = "MAF air flow rate"
Description[17] = "Throttle position"
Description[18] = "Commanded secondary air status"
Description[19] = "Oxygen sensors present (in 2 banks)"
Description[20] = "Oxygen Sensor 1, A: Voltage, B: Short term fuel trim"
Description[21] = "Oxygen Sensor 2, A: Voltage, B: Short term fuel trim"
Description[22] = "Oxygen Sensor 3, A: Voltage, B: Short term fuel trim"
Description[23] = "Oxygen Sensor 4, A: Voltage, B: Short term fuel trim"
Description[24] = "Oxygen Sensor 5, A: Voltage, B: Short term fuel trim"
Description[25] = "Oxygen Sensor 6, A: Voltage, B: Short term fuel trim"
Description[26] = "Oxygen Sensor 7, A: Voltage, B: Short term fuel trim"
Description[27] = "Oxygen Sensor 8, A: Voltage, B: Short term fuel trim"
Description[28] = "OBD standards this vehicle conforms to"
Description[29] = "Oxygen sensors present (in 4 banks)"
Description[30] = "Auxiliary input status"
Description[31] = "Run time since engine start"
Description[32] = "PIDs supported [21 - 40]"
Description[33] = "Distance traveled with malfunction indicator lamp (MIL) on"
Description[34] = "Fuel Rail Pressure (relative to manifold vacuum)"
Description[35] = "Fuel Rail Gauge Pressure (diesel, or gasoline direct injection)"
Description[36] = "Oxygen Sensor 1, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[37] = "Oxygen Sensor 2, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[38] = "Oxygen Sensor 3, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[39] = "Oxygen Sensor 4, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[40] = "Oxygen Sensor 5, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[41] = "Oxygen Sensor 6, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[42] = "Oxygen Sensor 7, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[43] = "Oxygen Sensor 8, AB: Fuel–Air Equivalence Ratio, CD: Voltage"
Description[44] = "Commanded EGR"
Description[45] = "EGR Error"
Description[46] = "Commanded evaporative purge"
Description[47] = "Fuel Tank Level Input"
Description[48] = "Warm-ups since codes cleared"
Description[49] = "Distance traveled since codes cleared"
Description[50] = "Evap. System Vapor Pressure"
Description[51] = "Absolute Barometric Pressure"
Description[52] = "Oxygen Sensor 1, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[53] = "Oxygen Sensor 2, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[54] = "Oxygen Sensor 3, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[55] = "Oxygen Sensor 4, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[56] = "Oxygen Sensor 5, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[57] = "Oxygen Sensor 6, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[58] = "Oxygen Sensor 7, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[59] = "Oxygen Sensor 8, AB: Fuel–Air Equivalence Ratio, CD: Current"
Description[60] = "Catalyst Temperature: Bank 1, Sensor 1"
Description[61] = "Catalyst Temperature: Bank 2, Sensor 1"
Description[62] = "Catalyst Temperature: Bank 1, Sensor 2"
Description[63] = "Catalyst Temperature: Bank 2, Sensor 2"
Description[64] = "PIDs supported [41 - 60]"
Description[65] = "Monitor status this drive cycle"
Description[66] = "Control module voltage"
Description[67] = "Absolute load value"
Description[68] = "Fuel–Air commanded equivalence ratio"
Description[69] = "Relative throttle position"
Description[70] = "Ambient air temperature"
Description[71] = "Absolute throttle position B"
Description[72] = "Absolute throttle position C"
Description[73] = "Absolute throttle position D"
Description[74] = "Absolute throttle position E"
Description[75] = "Absolute throttle position F"
Description[76] = "Commanded throttle actuator"
Description[77] = "Time run with MIL on"
Description[78] = "Time since trouble codes cleared"
Description[79] = "Maximum value for Fuel–Air equivalence ratio, oxygen sensor voltage, oxygen sensor current, and intake manifold absolute pressure"
Description[80] = "Maximum value for air flow rate from mass air flow sensor"
Description[81] = "Fuel Type"
Description[82] = "Ethanol fuel %"
Description[83] = "Absolute Evap system Vapor Pressure"
Description[84] = "Evap system vapor pressure"
Description[85] = "Short term secondary oxygen sensor trim, A: bank 1, B: bank 3"
Description[86] = "Long term secondary oxygen sensor trim, A: bank 1, B: bank 3"
Description[87] = "Short term secondary oxygen sensor trim, A: bank 2, B: bank 4"
Description[88] = "Long term secondary oxygen sensor trim, A: bank 2, B: bank 4"
Description[89] = "Fuel rail absolute pressure"
Description[90] = "Relative accelerator pedal position"
Description[91] = "Hybrid battery pack remaining life"
Description[92] = "Engine oil temperature"
Description[93] = "Fuel injection timing"
Description[94] = "Engine fuel rate"
Description[95] = "Emission requirements to which vehicle is designed"
Description[96] = "PIDs supported [61 - 80]"
Description[97] = "Driver's demand engine - percent torque"
Description[98] = "Actual engine - percent torque"
Description[99] = "Engine reference torque"

OBD_PID_COUNT = 0
OBD_DATA_DICT = {}
CRANE_SORTED_SENSOR_LIST = []
CRANE_SENSOR_LENGTH_DICT = {}
CURRENT_VERSION = ""
SPARE_LENGTH = 10

def getDataFromInfluxDB(INFLUXDB, CraneFullName):
    global OBD_PID_COUNT
    global OBD_DATA_DICT

    global CRANE_SORTED_SENSOR_LIST
    global CRANE_SENSOR_LENGTH_DICT

    global CURRENT_VERSION

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

    query1 = 'select * from ' + CraneFullName + '_SensorList'
    scan_data = client.query(query1)

    OBD_PID_COUNT = 0
    OBD_DATA_DICT = {}
    CRANE_SORTED_SENSOR_LIST  = []
    CRANE_SENSOR_LENGTH_DICT = {}


    sensor_dict = {}
    for value in scan_data.get_points() :
        if value['Type'] == 'OBD' :
            # get OBD_PID_COUNT
            # OBD 개수 카운팅
            OBD_PID_COUNT = OBD_PID_COUNT + 1

            # get OBD_DATA_DICT
            # key에 해당하는 센서이름 얻음
            hex_pid = int(value['Num'], 16)
            #hex_pid = value['Num']
            #OBD_DATA_DICT[hex_pdi] = key

        elif value['Type'] == 'Crane' :
            # get CRANE_SORTED_SENSOR_LIST
            # 정렬된 Crane Sensor List 
            #sensor_name = key
            sensor_name = value['SensorName']
            sort_num = int(value['Num'])
            sensor_dict[sort_num] = sensor_name

            # get CRANE_SENSOR_LENGTH_DICT
            # Sensor를 표기하는 범위 길이 구함
            max = str(int(value['RangeMax']))
            min = str(int(value['RangeMin']))
            decimal_places = value['DecimalPlaces']

            max_length = len(max)
            min_length = len(min)
            

            if max_length >= min_length :
                length = max_length
            else :
                length = min_length

            if float(min) < 0 :
                if min_length - 1 >= max_length :
                    length = min_length
                else :
                    length = max_length + 1

            if int(decimal_places) != 0 :
                length = length + 1 + int(decimal_places)

            CRANE_SENSOR_LENGTH_DICT[sensor_name] = length


    for key, value in sorted(sensor_dict.items()) :
        CRANE_SORTED_SENSOR_LIST.append(value)
    print(CRANE_SORTED_SENSOR_LIST)
  

 
    

def isNumber(s):
  try:
    float(s)
    return True
  except ValueError:
    return False

# def check_obd_pid(msg) :
#     msg_split = msg.split('|')

#     pid_count = 0

#     for i in range(0, len(msg_split)) :
#         try :
#             time = datetime.datetime.strptime(msg_split[i], "%Y-%m-%d %H:%M:%S")
#             break
#         except ValueError :
#             pid_count = pid_count + 1

#     return pid_count



# def get_pid_count(NameNode, table_name, table_seperate) :

#     connection = happybase.Connection(NameNode, autoconnect=False)
#     connection.open()

#     table = connection.table(table_seperate + table_name)

#     scan_data = table.scan()

#     pid_count = 0
#     for key, value in scan_data :
#         if value[b'Info:Type'] == 'OBD' :
#             pid_count = pid_count + 1
#     connection.close()
  
#     return pid_count

def get_pid_count() :
    return OBD_PID_COUNT

def get_obd_description(hex_pid) :
    hex_pid = int(hex_pid, 16)

    return OBD_DATA_DICT[hex_pid]

def get_sorted_sensor_list():
    return CRANE_SORTED_SENSOR_LIST

def get_sensor_length_dict():
    return CRANE_SENSOR_LENGTH_DICT

def get_hbase_sensor_data_len(sensor_length_dict) :
    sensor_data_len = 0

    for value in sensor_length_dict.values() :
        sensor_data_len = sensor_data_len + int(value)

    return sensor_data_len

# def get_obd_description(NameNode, hex_pid, table_seperate) :

#     connection = happybase.Connection(NameNode, autoconnect=False)
#     connection.open()

#     table_name = table_seperate + '_SensorList'
#     table = connection.table(table_name)
#     hex_pid = int(hex_pid, 16)

#     description = ""

#     for key, value in table.scan() :
#         # print(value[b'Info:Num'])
#         # hex_pid = int(hex_pid, 16)
#         value[b'Info:Num'] = int(value[b'Info:Num'],16)
#         if value[b'Info:Type'] == 'OBD' and value[b'Info:Num'] == hex_pid :
#             # print("=========if문===========")
#             description =  str(key)

#     connection.close()
#     return description

def hexFormat(hex_value):
    dec_value = int(hex_value, 16)

    hex_value = hex(dec_value).upper()
    hex_value = '0' * (2 - len(hex_value[2:])) + hex_value[2:]

    return hex_value

def decFormat(hex_value):
    dec_value = int(hex_value, 16)
    return dec_value


def OBD_isVailable(OBD_DataList):
    
    if len(OBD_DataList) != get_pid_count():
        print("OBD Data is invaild : Generator Data Length != HBase Data Length")
        print("Recv : "+ str(len(OBD_DataList)))
        print("Defined_HBase : " + str(get_pid_count()))
        return False

    for i in range(len(OBD_DataList)):
        obd_data = eval(OBD_DataList[i])
        hex_pid = obd_data[3]
        dec_pid = decFormat(hex_pid)     

        if dec_pid not in OBD_DATA_DICT.keys():
            print("OBD Data is invaild : pid(hex) '"+ str(hex_pid)+"' Not Defined in HBase")
            return False

    return True

def OBD_II_Parser(msg):

    msg = eval(msg)
   
    pid = int(msg[3], 16) # 16 -> 10

    hex_pid = hex(pid).upper() # 10 -> 16 Upper
    hex_pid = '0' * (2 - len(hex_pid[2:])) + hex_pid[2:]

    A = int(msg[4], 16)
    B = int(msg[5], 16)
    C = int(msg[6], 16)
    D = int(msg[7], 16)

    # print("PID : " + str(pid) + ", PID(Hex) : " + str(hex_pid) + ", Descripition : " + Description[pid] + ", A : " + str(A) + ", B : " + str(B) + ", C : " + str(C) + ", D : " + str(D))

    description = str(get_obd_description(str(hex_pid)))
    # print("=========test=========")
    # print(description)

    return description, pid, A, B, C, D

def OBD_II_Calc(pid, A, B, C, D):
    value = ""
    if pid == 0:
        temp = bin(A)
        temp = '0' * (8 - len(temp[2:])) + temp[2:]

        temp2 = bin(B)
        temp2 = '0' * (8 - len(temp2[2:])) + temp2[2:]

        temp3 = bin(C)
        temp3 = '0' * (8 - len(temp3[2:])) + temp3[2:]

        temp4 = bin(D)
        temp4 = '0' * (8 - len(temp4[2:])) + temp4[2:]

        sum_temp = temp + temp2 + temp3 + temp4

        for i in range(len(sum_temp)):
            if sum_temp[i] == '1':
                hex_val = hex(i+1)
                hex_val = '0' * (2 - len(hex_val[2:])) + hex_val[2:]
                value = value + hex_val + ", "

        value = value[:-2]

    elif pid == 1:
        pass

    elif pid == 2:
        pass

    elif pid == 3:
        pass

    elif pid == 4:
        value = round(100/255.0 * A, 2)
        

    elif pid == 5:
        value = round(A - 40, 2)

    elif pid == 6:
        value = round(100/128.0 * A - 100, 2)

    elif pid == 7:
        value = round(100/128.0 * A - 100, 2)

    elif pid == 8:
        value = round(100/128.0 * A - 100, 2)

    elif pid == 9:
        value = round(100/128.0 * A - 100, 2)

    elif pid == 10:
        value = round(3 * A, 2)

    elif pid == 11 :
        value = round(A, 2)

    elif pid == 12:
        value = round((256 * A + B) / 4.0, 2)

    elif pid == 13 :
        value = round(A, 2)

    elif pid == 14:
        value = round(A / 2.0 - 64, 2)

    elif pid == 15:
        value = round(A - 40, 2)

    elif pid == 16:
        value = round((256 * A + B) / 100.0, 2)

    elif pid == 17:
        value = round(100 / 255.0 * A, 2)

    elif pid == 18:
        pass
    elif pid == 19:
        pass
    elif pid == 20:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 21:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 22:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 23:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 24:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 25:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 26:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 27:
        value = str(round(A / 200.0, 2)) + "," + str(round(100/128.0 * B - 100, 2))

    elif pid == 28:
        pass
    elif pid == 29:
        pass
    elif pid == 30:
        pass
    elif pid == 31:
        value = round(256 * A + B, 2)

    elif pid == 32:
        temp = bin(A)
        temp = '0' * (8 - len(temp[2:])) + temp[2:]

        temp2 = bin(B)
        temp2 = '0' * (8 - len(temp2[2:])) + temp2[2:]

        temp3 = bin(C)
        temp3 = '0' * (8 - len(temp3[2:])) + temp3[2:]

        temp4 = bin(D)
        temp4 = '0' * (8 - len(temp4[2:])) + temp4[2:]

        sum_temp = temp + temp2 + temp3 + temp4

        for i in range(len(sum_temp)):
            if sum_temp[i] == '1':
                hex_val = hex(i + 33)
                hex_val = '0' * (2 - len(hex_val[2:])) + hex_val[2:]
                value = value + hex_val + ", "

        value = value[:-2]
    elif pid == 33:
        value = round(256 * A + B, 2)

    elif pid == 34:
        value = round(0.079 * (256 * A + B), 2)

    elif pid == 35:
        value = round(10 * (256 * A + B), 2)

    elif pid == 36:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 37:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 38:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 39:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 40:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 41:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 42:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 43:
        value = str(round(2 / 65536.0 * (256 * A + B), 2)) + "," + str(round(8 / 65536.0 * (256 * C + D), 2))

    elif pid == 44:
        value = round(100/255.0 * A, 2)

    elif pid == 45:
        value = round(100/128.0 * A - 100, 2)

    elif pid == 46:
        value = round(100/255.0 * A, 2)

    elif pid == 47:
        value = round(100 / 255.0 * A, 2)

    elif pid == 48:
        value = round(A, 2)

    elif pid == 49:
        value = round(256 * A + B, 2)

    elif pid == 50:
        value = round((256 * A + B) / 4.0, 2)

    elif pid == 51:
        value = round(A, 2)

    elif pid == 52:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 53:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 54:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 55:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 56:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 57:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 58:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 59:
        value = str(round(2 / 65536.0 * ((256 * A + B)), 2)) + "," +  str(round(C + D/256.0 - 128, 2))

    elif pid == 60:
        value = round((256 * A + B) / 10.0 - 40, 2)

    elif pid == 61:
        value = round((256 * A + B) / 10.0 - 40, 2)

    elif pid == 62:
        value = round((256 * A + B) / 10.0 - 40, 2)

    elif pid == 63:
        value = round((256 * A + B) / 10.0 - 40, 2)

    elif pid == 64:
        temp = bin(A)
        temp = '0' * (8 - len(temp[2:])) + temp[2:]

        temp2 = bin(B)
        temp2 = '0' * (8 - len(temp2[2:])) + temp2[2:]

        temp3 = bin(C)
        temp3 = '0' * (8 - len(temp3[2:])) + temp3[2:]

        temp4 = bin(D)
        temp4 = '0' * (8 - len(temp4[2:])) + temp4[2:]

        sum_temp = temp + temp2 + temp3 + temp4

        for i in range(len(sum_temp)):
            if sum_temp[i] == '1':
                hex_val = hex(i + 65)
                hex_val = '0' * (2 - len(hex_val[2:])) + hex_val[2:]
                value = value + hex_val + ", "

        value = value[:-2]

    elif pid == 65:
        pass

    elif pid == 66:
        value = round((256 * A + B) / 1000.0, 2)

    elif pid == 67:
        value = round(100/255.0 * (256 * A + B), 2)

    elif pid == 68:
        value = round(2/65536.0 * (256 * A + B) ,2)

    elif pid == 69:
        value = round(100 / 255.0 * A ,2)

    elif pid == 70:
        value = round(A - 40, 2)

    elif pid == 71:
        value = round(100 / 255.0 * A, 2)

    elif pid == 72:
        value = round(100 / 255.0 * A, 2)

    elif pid == 73:
        value = round(100 / 255.0 * A, 2)

    elif pid == 74:
        value = round(100 / 255.0 * A, 2)

    elif pid == 75:
        value = round(100 / 255.0 * A, 2)

    elif pid == 76:
        value = round(100 / 255.0 * A, 2)

    elif pid == 77:
        value = round(256 * A + B, 2)

    elif pid == 78:
        value = round(256 * A + B, 2)

    elif pid == 79:     
        value = str(round(A, 2)) + "," +  str(round(B, 2)) + "," + str(round(C, 2)) + "," + str(round(D*10, 2))

    elif pid == 80:
        pass

    elif pid == 81:
        pass

    elif pid == 82:
        value = round(100 / 255.0 * A, 2)

    elif pid == 83:
        value = round((256 * A + B) / 200.0, 2)

    elif pid == 84:
        value = round(((A * 256) + B) - 32767, 2)

    elif pid == 85:
        value = str(round(100 / 128.0 * A - 100, 2)) + "," + str(round(100 / 128.0 * B - 100, 2))

    elif pid == 86:
        value =  str(round(100 / 128.0 * A - 100, 2)) + "," + str(round(100 / 128.0 * B - 100, 2))

    elif pid == 87:
        value =  str(round(100 / 128.0 * A - 100, 2)) + "," + str(round(100 / 128.0 * B - 100, 2))

    elif pid == 88:
        value =  str(round(100 / 128.0 * A - 100, 2)) + "," + str(round(100 / 128.0 * B - 100, 2))

    elif pid == 89:
        value = round(10 * (256 * A + B), 2)

    elif pid == 90:
        value = round(100 / 255.0 * A, 2)

    elif pid == 91:
        value = round(100 / 255.0 * A, 2)

    elif pid == 92:
        value = round(A - 40, 2)

    elif pid == 93:
        value = round((256 * A + B) / 128.0 - 210, 2)

    elif pid == 94:
        value = round((256 * A + B) / 20.0, 2)

    elif pid == 95:
        pass

    elif pid == 96:
        temp = bin(A)
        temp = '0' * (8 - len(temp[2:])) + temp[2:]

        temp2 = bin(B)
        temp2 = '0' * (8 - len(temp2[2:])) + temp2[2:]

        temp3 = bin(C)
        temp3 = '0' * (8 - len(temp3[2:])) + temp3[2:]

        temp4 = bin(D)
        temp4 = '0' * (8 - len(temp4[2:])) + temp4[2:]

        sum_temp = temp + temp2 + temp3 + temp4

        for i in range(len(sum_temp)):
            if sum_temp[i] == '1':
                hex_val = hex(i + 97)
                hex_val = '0' * (2 - len(hex_val[2:])) + hex_val[2:]
                value = value + hex_val + ", "

        value = value[:-2]

    elif pid == 97:
        value = round(A - 125, 2)

    elif pid == 98:
        value = round(A - 125, 2)

    elif pid == 99:
        value = round(256 * A + B, 2)



    print("Calculation => Value : " + str(value))
    print
    return value

# parsing 순서대로 센서를 정렬해서 리스트 반환
# def get_sorted_sensor_list(NameNode, table_name, table_seperate) :
    
#     connection = happybase.Connection(NameNode, autoconnect=False)
#     connection.open()

#     table = connection.table(table_seperate + table_name)
#     sensor_dict = {}
#     # print(type(sensor_dict))

#     for key, value in table.scan() :
#         if value[b'Info:Type'] != 'OBD' :
#             sensor_name = key
#             sort_num = int(value[b'Info:Num'])
#             sensor_dict[sort_num] = sensor_name

#     sorted_sensor_list = []

#     for key, value in sorted(sensor_dict.items()) :
#         sorted_sensor_list.append(value)
   
#     connection.close()
#     return sorted_sensor_list

# # 센서 이름을 Key, Parsing을 위한 길이 값을 Value로 가지는 딕셔너리 반환
# def get_sensor_length_dict(NameNode, table_name, table_seperate) :
   
#     connection = happybase.Connection(NameNode, autoconnect=False)
#     connection.open()

#     table = connection.table(table_seperate + table_name)

#     sensor_length_dict = {}

#     for key, value in table.scan() :
#         if value[b'Info:Type'] != 'OBD' :
#             sensor_name = key
#             max = value[b'Range:Max']
#             min = value[b'Range:Min']
#             decimal_places = value[b'Info:DecimalPlaces']

#             max_length = len(max)
#             min_length = len(min)

#             if max_length >= min_length :
#                 length = max_length
#             else :
#                 length = min_length

#             if float(min) < 0 :
#                 if min_length - 1 >= max_length :
#                     length = min_length
#                 else :
#                     length = max_length + 1

#             if int(decimal_places) != 0 :
#                 length = length + 1 + int(decimal_places)

#             sensor_length_dict[sensor_name] = length
#     # print("======================")
#     # print(sensor_length_dict)
#     print
#     connection.close()
#     return sensor_length_dict


# def get_hbase_sensor_data_len(sensor_length_dict) :
#     sensor_data_len = 0

#     for value in sensor_length_dict.values() :
#         sensor_data_len = sensor_data_len + int(value)

#     return sensor_data_len

# 숫자 판별 함수
def isNumber(s):
    try:
        float(s)
        return True

    except ValueError:
        return False


#센서 데이터를 파싱하여 센서 이름을 Key, 센서 값을 Value로 가지는 딕셔너리 반환 
def sensor_data_parser(sensor_data) :
    # 인풋 데이터
    sensor_data = str(sensor_data)

    # HBase 데이터
    sorted_sensor_list = get_sorted_sensor_list()
    sensor_length_dict = get_sensor_length_dict()

    sensor_data_len = get_hbase_sensor_data_len(sensor_length_dict)
    

    #STX = sensor_data[0:1]
    #sensor_data =sensor_data[0:]
    #ETX = sensor_data[-1]
    #seperator = sensor_data[-2:-1]
    #sensor_data = sensor_data[ : -2]

    sensor_data = sensor_data[1:]

    if len(sensor_data) != (sensor_data_len + SPARE_LENGTH) :
        print("Crane Data is invaild : Generator Data Length != HBase Data Length")
        print("Recv : " + str(len(sensor_data)))
        print("Defined_HBase : "+ str(sensor_data_len +SPARE_LENGTH))
        return 'ignore'

    sensor_data_list = []
    sensor_dict = {}

    for sensor_name in sorted_sensor_list :

        individual_sensor_value = sensor_data[0 : sensor_length_dict[sensor_name]]
        print(individual_sensor_value)
        if not isNumber(individual_sensor_value):
            print("Crane Data is invaild : Generator Data != HBase Data")
            return 'ignore'

        sensor_data_list.append(individual_sensor_value)
        sensor_dict[sensor_name] = individual_sensor_value

        sensor_data = sensor_data[sensor_length_dict[sensor_name] : ]


    print("=> Parsing data : " + str(sensor_data_list))
    print

    return sensor_dict


def crane_data_parser(crane_data_encoded_hex):
    crane_data = binascii.unhexlify(crane_data_encoded_hex)
    crane_data = str(crane_data)
 
    crane_data_list = ["" for i in range(31)]

    STX = crane_data[0:1]

    crane_data_list[0] = crane_data[1:7] # main_weight
    crane_data_list[1] = crane_data[7:12] # main_ad
    crane_data_list[2]  = crane_data[12:17] # aux_weight
    crane_data_list[3]  = crane_data[17:22]  # aux_ad
    crane_data_list[4]  = crane_data[22:27]  # boom_angle
    crane_data_list[5]  = crane_data[27:32]  # boom_ad
    crane_data_list[6]  = crane_data[32:36]  # working_radius
    crane_data_list[7]  = crane_data[36:41]  # limit
    crane_data_list[8]  = crane_data[41:46]  # boom_lenght
    crane_data_list[9]  = crane_data[46:47]  # hoist
    crane_data_list[10]  = crane_data[47:49]  # main_fall
    crane_data_list[11]  = crane_data[49:53]  # adc0_12V
    crane_data_list[12]  = crane_data[53:57]  # adc1_5V
    crane_data_list[13]  = crane_data[57:61]  # adc2_EXC
    crane_data_list[14]  = crane_data[61:65]  # adc3_EXC
    crane_data_list[15]  = crane_data[65:69]  # adc4_12V
    crane_data_list[16]  = crane_data[69:73]  # adc5_5V
    crane_data_list[17]  = crane_data[73:79]  # plus_12V
    crane_data_list[18]  = crane_data[79:84]  # plus_5V
    crane_data_list[19]  = crane_data[84:89]  # plus_EXC
    crane_data_list[20]  = crane_data[89:94]  # minus_EXC
    crane_data_list[21]  = crane_data[94:100]  # minus_12V
    crane_data_list[22]  = crane_data[100:105]  # minus_5V
    crane_data_list[23]  = crane_data[105:107]  # error_code
    crane_data_list[24]  = crane_data[107:108]  # block_led
    crane_data_list[25]  = crane_data[108:109]  # ign_block
    crane_data_list[26]  = crane_data[109:110]  # buzzer_led
    crane_data_list[27]  = crane_data[110:111]  # A2B_led
    crane_data_list[28]  = crane_data[111:112]  # max_led
    crane_data_list[29]  = crane_data[112:113]  # min_led
    crane_data_list[30]  = crane_data[113:115]  # aux_fall_num


    seperator = crane_data[115:116]  # serperator
    ETX = crane_data[116:117]

    print("--------Crane Data Parsing--------")
    print("=> RECV data : " + crane_data_encoded_hex)
    print
    print("=> Parsing data : " + str(crane_data_list))
    print


    return crane_data_list


def make_dict(Time, CraneFullName, OBD_II_data_dict, sensor_data_dict) :
    data_dict = {}
    if len(OBD_II_data_dict) != 0:
        data_dict.update(OBD_II_data_dict)
    if len(sensor_data_dict) != 0:
        data_dict.update(sensor_data_dict)

    data_dict['Time'] = Time
    data_dict['CraneFullName'] = CraneFullName
    #data_dict['Version'] = version

    return data_dict


def makeJson(data_dict) :

    jsonString = json.dumps(data_dict)

    return jsonString








