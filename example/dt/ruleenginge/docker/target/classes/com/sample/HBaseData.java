package com.sample;

//===========hadoop&hbase================
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.hbase.client.HBaseAdmin;
import org.apache.hadoop.hbase.TableName;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.HTableDescriptor;
import org.apache.hadoop.hbase.HColumnDescriptor;
import org.apache.hadoop.hbase.mapreduce.TableInputFormat;
import org.apache.hadoop.hbase.KeyValue.Type;
import org.apache.hadoop.hbase.HConstants;
import org.apache.hadoop.hbase.util.Bytes;
import org.apache.hadoop.hbase.CellUtil;
import org.apache.hadoop.hbase.Cell;
import org.apache.hadoop.hbase.client.Table;
import org.apache.hadoop.hbase.client.Put;
import org.apache.hadoop.hbase.client.Get;
import org.apache.hadoop.hbase.client.Scan;
import org.apache.hadoop.hbase.client.Result;
import org.apache.hadoop.hbase.client.ResultScanner;

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

public class HBaseData{
    // =============== table list======================
    private final String data_table_tail = "_Data";
    private final String result_table_tail = "_Result";
    private final String errorList_table_tail = "_ErrorList";
    private final String sensorList_table_tail = "_SensorList";
    private final String version_table_tail = "_Version";

    // ==============table columns ==========================
    private final List<String> data_cf_list = Arrays.asList("Value");
    private final List<String> result_cf_list = Arrays.asList("Result", "ResultErrorCode", "Normal_Min", "Normal_Max");



    org.apache.hadoop.conf.Configuration hconf;
    Configure c;

    HashMap<String, String> error_list_map;
    HashMap<String, String> range_map;

    public HBaseData(){
        c = Configure.getInstance();
        this.hconf = createHBaseConfigure();

    }
    // Hbase 설정 함수
    private org.apache.hadoop.conf.Configuration createHBaseConfigure() {
        org.apache.hadoop.conf.Configuration hconf = HBaseConfiguration.create();
        String zookeeper_quorum = c.ZOOKEEPER_QUORUM;
        hconf.set("hbase.zookeeper.property.clientPort",c.ZOOKEEPER_PORT);
        hconf.set("hbase.zookeeper.quorum", zookeeper_quorum);
        hconf.set("hbase.cluster.distributed","true");
        return hconf;
    }

    // Hbase의 ErrorList 테이블로 부터 데이터를 읽음
    public HashMap<String, String> getErrorListMap(String errorList_table_tail) throws IOException {
      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(this.hconf);
      Table table = connection.getTable(TableName.valueOf(errorList_table_tail));

      HashMap <String, String> states = new HashMap<String, String>();
      ResultScanner scan = table.getScanner(new Scan());

      for (Result r : scan) {
        String errorBitNum = Bytes.toString(r.getRow());
       
        for (Cell kv : r.listCells()) {
          String errorName = Bytes.toString(CellUtil.cloneValue(kv));
          states.put(errorBitNum, errorName);
        }
      }
     
      connection.close();
      return states;
    } 


    // Hbase SensorList 테이블로 부터 각 센서별 Range 범위 데이터를 얻어옴
    public HashMap<String, String> getRangeMap( String sensorList_table_tail) throws IOException {
      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(this.hconf);
      Table range_table = connection.getTable(TableName.valueOf(sensorList_table_tail));

      HashMap<String, String> states = new HashMap<String, String>();
      ResultScanner scan = range_table.getScanner(new Scan());

      String range_max = "";
      String range_min = "";

      for ( Result r : scan ) {
        String sensorName = Bytes.toString(r.getRow());

        for ( Cell kv : r.listCells()) {
          String col_family = Bytes.toString(CellUtil.cloneFamily(kv));

          if (col_family.equals("NormalRange")) {

            String col_name =  Bytes.toString(CellUtil.cloneQualifier(kv));
         
            if (col_name.equals("Max")) {
                  range_max = Bytes.toString(CellUtil.cloneValue(kv));
                  states.put(sensorName + "_normal_max", range_max);
            }
            else {
                  range_min = Bytes.toString(CellUtil.cloneValue(kv));
                  states.put(sensorName + "_normal_min", range_min);
            }
           
          }
        }
      }
      connection.close();
      return states;
    }
    public void sensor_error_load(String CraneFullName)  throws IOException {
        error_list_map = getErrorListMap(CraneFullName + errorList_table_tail);
        range_map = getRangeMap(CraneFullName + sensorList_table_tail);
    }

    // CEP Engine 분석 결과를 Hbase에 저장
    public void saveHBaseResult(Sensor obj_sensor) throws IOException {
        
      String rowKey = obj_sensor.get_dataMap().get("Time");
      String CraneFullName = obj_sensor.get_dataMap().get("CraneFullName");

      org.apache.hadoop.conf.Configuration hconf = createHBaseConfigure();
      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(hconf);
      Table table = connection.getTable(TableName.valueOf(CraneFullName + result_table_tail));
      Put result_put = new Put(Bytes.toBytes(rowKey));

      // HashMap<String, String> error_list_map = getErrorListMap(CraneFullName + errorList_table_tail);
      // HashMap<String, String> range_map = getRangeMap(CraneFullName + sensorList_table_tail);


      for(HashMap.Entry<String, String> entry : obj_sensor.resultMap.entrySet()) {
        result_put.add(Bytes.toBytes(result_cf_list.get(0)), Bytes.toBytes(entry.getKey()),  Bytes.toBytes(entry.getValue()));
      }

      int count = 0;

      result_put.add(Bytes.toBytes(result_cf_list.get(1)), Bytes.toBytes("RuleErrorCode"),  Bytes.toBytes(obj_sensor.get_ruleErrorCode()));

      for(HashMap.Entry<String, String> entry : error_list_map.entrySet()) {
        if((obj_sensor.error_value & (int)(Math.pow(2, Integer.parseInt(entry.getKey()) - 1))) != 0){ // set bit
                     count = count + 1;
                     result_put.add(Bytes.toBytes(result_cf_list.get(1)), Bytes.toBytes("detail_"+ Integer.toString(count)),  Bytes.toBytes(entry.getValue()));
        }
      }

      for(HashMap.Entry<String, String> entry : obj_sensor.resultMap.entrySet()) {
          try {
            result_put.add(Bytes.toBytes(result_cf_list.get(2)), Bytes.toBytes(entry.getKey()),  Bytes.toBytes(range_map.get(entry.getKey()+"_normal_min")));
            result_put.add(Bytes.toBytes(result_cf_list.get(3)), Bytes.toBytes(entry.getKey()),  Bytes.toBytes(range_map.get(entry.getKey()+"_normal_max")));  
          } catch (Exception e) {
            e.printStackTrace();
          }

      }

      table.put(result_put);
      connection.close();

    }


    // Hbase SensorList 테이블에서 센서 리스트를 얻어옴
    public List<String> getSensorList(String sensorList_table_tail) throws IOException {
      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(this.hconf);
      Table sensorList_table = connection.getTable(TableName.valueOf(sensorList_table_tail));
      List<String> states = new ArrayList<String>();
      ResultScanner scan = sensorList_table.getScanner(new Scan());
      for (Result r : scan) {
        String key = Bytes.toString(r.getRow());
        states.add(key);
      }
      connection.close();
      return states;
    }
    // Hbase Version 테이블에서 현재 버전을 얻어옴
    public String getVersion(String CraneFullName) throws IOException {
      String version_tableName = CraneFullName + version_table_tail;
      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(this.hconf);
      Table version_table = connection.getTable(TableName.valueOf(version_tableName));
      String result_version = "";
      ResultScanner scan = version_table.getScanner(new Scan());
      for (Result r : scan) {
        String key = Bytes.toString(r.getRow());
        result_version = Bytes.toString(r.getValue(Bytes.toBytes("Version"), Bytes.toBytes("Version")));
        
      }
      connection.close();
      return result_version;
    }
    
}