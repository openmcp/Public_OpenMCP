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


public class HBaseReader {

    private final String node_table_name = "raw_data_node";
    private final String pod_table_name = "raw_data_pod";
    private final String zookeeper_quorum = "zk-cs.datacenter.svc.cluster.local";
    private final String zookeeper_port = "2181";;

    //private final String zookeeper_quorum = "10.0.4.103";
    //private final String zookeeper_port = "32115";
    
    org.apache.hadoop.conf.Configuration hconf;

    public HBaseReader() {
        this.hconf = createHBaseConfigure();
    }
    // Hbase 설정 함수
    private org.apache.hadoop.conf.Configuration createHBaseConfigure() {
        org.apache.hadoop.conf.Configuration hconf = HBaseConfiguration.create();
        hconf.set("hbase.zookeeper.property.clientPort", zookeeper_port);
        hconf.set("hbase.zookeeper.quorum", zookeeper_quorum);
        hconf.set("hbase.cluster.distributed","true");
        return hconf;
    }
    // Hbase SensorList 테이블로 부터 각 센서별 Range 범위 데이터를 얻어옴
    public HashMap<String, HashMap<String,HashMap<String, String>>> getRawDataMap(String table_name, String start_time, String last_time) throws IOException {
      System.out.println("  [Table Name] " + table_name);


      org.apache.hadoop.hbase.client.Connection connection = org.apache.hadoop.hbase.client.ConnectionFactory.createConnection(this.hconf);

      Table raw_table = connection.getTable(TableName.valueOf(table_name));

      HashMap<String, HashMap<String, HashMap<String, String>>> data_map = new HashMap<String, HashMap<String,HashMap<String, String>>>();

      Scan scan = new Scan(Bytes.toBytes(start_time), Bytes.toBytes(last_time));
      ResultScanner scanner = raw_table.getScanner(scan);

      for ( Result r : scanner ) {
        String Time = Bytes.toString(r.getRow());
        System.out.println("    [FindData] "+ Time);
        
        for ( Cell kv : r.listCells()) {
          String col_family = Bytes.toString(CellUtil.cloneFamily(kv));
          String col_name = Bytes.toString(CellUtil.cloneQualifier(kv));
          String value = Bytes.toString(CellUtil.cloneValue(kv));

          // Node or Pod Name
          // Time
          // Resource ( CPU, Mem ... )
          // Value 
          data_map = put(data_map, Time, col_name, col_family, value);
        }
      }

      connection.close();
      return data_map;
    }
    public HashMap<String, HashMap<String, HashMap<String, String>>> put(HashMap<String, HashMap<String, HashMap<String, String>>> data_map, String time, String name, String resource, String value){
        if(data_map.get(time) == null) { 
            data_map.put(time, new HashMap<String, HashMap<String, String>>());
        }
        if(data_map.get(time).get(name) == null){
            data_map.get(time).put(name, new HashMap<String, String>());
        }
        data_map.get(time).get(name).put(resource, value);
        return data_map;
    }
}



