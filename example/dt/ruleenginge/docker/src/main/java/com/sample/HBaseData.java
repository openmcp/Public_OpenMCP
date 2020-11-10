package com.sample;
// ================InfluxDB===============
import org.influxdb.InfluxDB;
import org.influxdb.InfluxDBFactory;
import org.influxdb.dto.Point;
import org.influxdb.dto.Point.Builder;
import org.influxdb.dto.BatchPoints;
import org.influxdb.dto.QueryResult;
import org.influxdb.dto.Query;
//===========java===============
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.IOException;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Arrays;
import java.util.concurrent.TimeoutException;
import java.nio.file.Paths;
import java.nio.file.Files;
import java.util.concurrent.TimeUnit;

public class HBaseData{
    // =============== table list======================
    private final String result_table_tail = "_Result";
    private final String version_table_tail = "_Version";
    private final String sensorList_table_tail = "_SensorList";
    private final String errorList_table_tail = "_ErrorList";


    InfluxDB influxDB;
    Configure c;

    HashMap<String, String> error_list_map;
    HashMap<String, String> range_map;

    public HBaseData(){
        c = Configure.getInstance();
        this.influxDB = InfluxDBFactory.connect("http://"+c.INFLUXDB+":"+c.INFLUX_PORT, c.INFLUX_USER, c.INFLUX_PW);
        influxDB.setDatabase("mydb");

    }
    public void sensor_error_load(String CraneFullName)  throws IOException {
        error_list_map = getErrorListMap(CraneFullName + errorList_table_tail);
        range_map = getRangeMap(CraneFullName + sensorList_table_tail);
    }

    public HashMap<String, String> getErrorListMap(String error_tableName) throws IOException {
      QueryResult queryResult = this.influxDB.query(new Query("SELECT * FROM "+error_tableName, "mydb"));

      HashMap <String, String> states = new HashMap<String, String>();
      List<List<Object>> objectList = queryResult.getResults().get(0).getSeries().get(0).getValues();
      for (int i = 0; i<objectList.size(); i++) {
        String errorBitNum = objectList.get(i).get(2).toString();
        String errorName = objectList.get(i).get(1).toString();
        states.put(errorBitNum, errorName);
      }


      return states;
    }
    // Hbase SensorList 테이블로 부터 각 센서별 Range 범위 데이터를 얻어옴
    public HashMap<String, String> getRangeMap( String sensorList_tableName) throws IOException {
      QueryResult queryResult = this.influxDB.query(new Query("SELECT * FROM "+sensorList_tableName, "mydb"));

      HashMap<String, String> states = new HashMap<String, String>();

      List<List<Object>> objectList = queryResult.getResults().get(0).getSeries().get(0).getValues();
      for (int i=0; i< objectList.size(); i++) {
        String NormalRangeMax = objectList.get(i).get(2).toString();
        String NormalRangeMin = objectList.get(i).get(3).toString();
        String SensorName = objectList.get(i).get(7).toString();

        states.put(SensorName + "_normal_max", NormalRangeMax);
        states.put(SensorName + "_normal_min", NormalRangeMin);

      }
      return states;
    }



    // CEP Engine 분석 결과를 Hbase에 저장
    public void saveHBaseResult(Sensor obj_sensor) throws IOException {
      String dbName = "mydb";
      String time_t = obj_sensor.get_dataMap().get("Time");
      String CraneFullName = obj_sensor.get_dataMap().get("CraneFullName");


      String measurementName = CraneFullName + result_table_tail;
      Builder builder = Point.measurement(measurementName).time(System.currentTimeMillis(), TimeUnit.MILLISECONDS);

      String field;
      String value;
      for(HashMap.Entry<String, String> entry : obj_sensor.resultMap.entrySet()) {
        field = entry.getKey();
        value = entry.getValue();
        builder = builder.addField(field, value);
      }

      field = "RuleErrorCode";
      value = obj_sensor.get_ruleErrorCode();
      builder = builder.addField(field, value);


      int count = 0;
      for(HashMap.Entry<String, String> entry : error_list_map.entrySet()) {
        if((obj_sensor.error_value & (int)(Math.pow(2, Integer.parseInt(entry.getKey()) - 1))) != 0){ // set bit
             count = count + 1;
             field = "detail_"+ Integer.toString(count);
             value = entry.getValue();
             builder = builder.addField(field, value);
        }
      }
      Point point = builder.build();

      BatchPoints batchPoints = BatchPoints.database(dbName).tag("async", "true").build();
      batchPoints.point(point);

      influxDB.write(batchPoints);
    }

    public String getVersion(String CraneFullName) throws IOException {
        String version_tableName = CraneFullName + version_table_tail;
        QueryResult queryResult = this.influxDB.query(new Query("SELECT Version FROM "+version_tableName, "mydb"));
        String result_version = queryResult.getResults().get(0).getSeries().get(0).getValues().get(0).get(1).toString();

        return result_version;
    }



    
}
