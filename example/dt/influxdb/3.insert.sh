IP=10.0.3.20
PORT=31211

printf "[DROP] Database\n"
curl -POST http://${IP}:${PORT}/query --data-urlencode "q=DROP DATABASE mydb"

printf "[Create] Database\n"
curl -POST http://${IP}:${PORT}/query --data-urlencode "q=CREATE DATABASE mydb"
 
printf "[Create] Retention Policy\n"
curl -POST http://${IP}:${PORT}/query --data-urlencode "q=CREATE RETENTION POLICY myrp ON mydb DURATION 365d REPLICATION 1 DEFAULT"

printf "[INSERT] SHINHAN_Crane_1_SensorList Data\n"

curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Weight,Type=Crane DecimalPlaces=1,Num=1,NormalRangeMax=90,NormalRangeMin=0,RangeMax=999,RangeMin=-999'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Length,Type=Crane DecimalPlaces=1,Num=2,NormalRangeMax=60,NormalRangeMin=10,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Angle,Type=Crane DecimalPlaces=1,Num=3,NormalRangeMax=85,NormalRangeMin=30,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Work_Radius,Type=Crane DecimalPlaces=1,Num=4,NormalRangeMax=50,NormalRangeMin=0,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Height,Type=Crane DecimalPlaces=1,Num=5,NormalRangeMax=100,NormalRangeMin=20,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Limit_Load,Type=Crane DecimalPlaces=1,Num=6,NormalRangeMax=90,NormalRangeMin=0,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Fall,Type=Crane DecimalPlaces=0,Num=7,NormalRangeMax=90,NormalRangeMin=2,RangeMax=99,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Error_Code,Type=Crane DecimalPlaces=0,Num=8,NormalRangeMax=99,NormalRangeMin=1,RangeMax=99,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Block,Type=Crane DecimalPlaces=0,Num=9,NormalRangeMax=1,NormalRangeMin=0,RangeMax=1,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Alert,Type=Crane DecimalPlaces=0,Num=10,NormalRangeMax=1,NormalRangeMin=0,RangeMax=1,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Upper_Bound,Type=Crane DecimalPlaces=0,Num=11,NormalRangeMax=1,NormalRangeMin=0,RangeMax=1,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Lower_Bound,Type=Crane DecimalPlaces=0,Num=12,NormalRangeMax=1,NormalRangeMin=0,RangeMax=1,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Turning,Type=Crane DecimalPlaces=1,Num=13,NormalRangeMax=360,NormalRangeMin=0,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_1,Type=Crane DecimalPlaces=1,Num=14,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_2,Type=Crane DecimalPlaces=1,Num=15,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_3,Type=Crane DecimalPlaces=1,Num=16,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_4,Type=Crane DecimalPlaces=1,Num=17,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_5,Type=Crane DecimalPlaces=1,Num=18,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_6,Type=Crane DecimalPlaces=1,Num=19,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_7,Type=Crane DecimalPlaces=1,Num=20,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_8,Type=Crane DecimalPlaces=1,Num=21,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_9,Type=Crane DecimalPlaces=1,Num=22,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_10,Type=Crane DecimalPlaces=1,Num=23,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Left_11,Type=Crane DecimalPlaces=1,Num=24,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_1,Type=Crane DecimalPlaces=1,Num=25,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_2,Type=Crane DecimalPlaces=1,Num=26,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_3,Type=Crane DecimalPlaces=1,Num=27,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_4,Type=Crane DecimalPlaces=1,Num=28,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_5,Type=Crane DecimalPlaces=1,Num=29,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_6,Type=Crane DecimalPlaces=1,Num=30,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_7,Type=Crane DecimalPlaces=1,Num=31,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_8,Type=Crane DecimalPlaces=1,Num=32,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_9,Type=Crane DecimalPlaces=1,Num=33,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_10,Type=Crane DecimalPlaces=1,Num=34,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Right_11,Type=Crane DecimalPlaces=1,Num=35,NormalRangeMax=90,NormalRangeMin=0,RangeMax=99,RangeMin=-99'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Latitude,Type=Crane DecimalPlaces=6,Num=36,NormalRangeMax=999,NormalRangeMin=0,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Longitude,Type=Crane DecimalPlaces=6,Num=37,NormalRangeMax=999,NormalRangeMin=0,RangeMax=999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Standard_Date,Type=Crane DecimalPlaces=0,Num=38,NormalRangeMax=991999,NormalRangeMin=0,RangeMax=999999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=Standard_Time,Type=Crane DecimalPlaces=0,Num=39,NormalRangeMax=125959,NormalRangeMin=0,RangeMax=999999,RangeMin=0'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_SensorList,SensorName=AML,Type=Crane DecimalPlaces=0,Num=40,NormalRangeMax=1,NormalRangeMin=0,RangeMax=1,RangeMin=0'

printf "sleeping for a second then inserting more data\n"
sleep 1

printf "[INSERT] SHINHAN_Crane_1_ErrorList Data\n"
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=1 ErrorName="Weight_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=2 ErrorName="Length_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=3 ErrorName="Limit_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=4 ErrorName="Angle_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=5 ErrorName="Height_Fault"'

curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=6 ErrorName="Limit_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=7 ErrorName="Fall_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=8 ErrorName="Alert_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=9 ErrorName="Turning_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=10 ErrorName="Upper_Fault"'

curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=11 ErrorName="Lower_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=12 ErrorName="Swing_Fault"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_ErrorList,ErrorNum=13 ErrorName="Left1_Fault"'

printf "sleeping for a second then inserting more data\n"
sleep 1

printf "[INSERT] SHINHAN_Crane_1_Rule Data\n"
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_Rule,RuleName=Rule1 ErrorName="Weight_Fault",Condition1="Weight >= 70"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_Rule,RuleName=Rule2 ErrorName="Length_Fault",Condition1="Length > 50 && Fall >= 45"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_Rule,RuleName=Rule3 ErrorName="Limit_Fault",Condition1="Limit_Load > 45"'
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_Rule,RuleName=Rule4 ErrorName="Angle_Fault",Condition1="Angle <= 30"'

printf "sleeping for a second then inserting more data\n"
sleep 1

printf "[INSERT] SHINHAN_Crane_1_Version Data\n"
curl -XPOST http://${IP}:${PORT}/write?db=mydb --data-binary 'SHINHAN_Crane_1_Version,VersionName=Version Version="1"'

#printf "\nSHOW MEASUREMENTS\n"
#curl -G http://${IP}:${PORT}/query --data-urlencode "db=mydb" --data-urlencode "q=SHOW MEASUREMENTS" --data-urlencode "pretty=true"
#
#printf "\nSHOW SERIES\n"
#curl -G http://${IP}:${PORT}/query --data-urlencode "db=mydb" --data-urlencode "q=SHOW SERIES" --data-urlencode "pretty=true"
#
#printf "\nSHOW TAG VALUES WITH KEY = service\n"
#curl -G http://${IP}:${PORT}/query --data-urlencode "db=mydb" --data-urlencode "q=SHOW TAG VALUES WITH KEY = service" --data-urlencode "pretty=true"
