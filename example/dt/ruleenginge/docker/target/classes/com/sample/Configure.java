package com.sample;

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

//===========json convert============
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;

public class Configure{
    // ==============configure file 내용 저장 변수================
    public static String QUEUE_NAME = "";
    public static String QUEUE_TOPIC = "";
    public static String ZOOKEEPER_PORT = "";
    public static String ZOOKEEPER_QUORUM = "";
    public static String MKDRL_PATH = "";
    public static String RULE_PATH="";
    public static String CONFIG_PATH="";

    private static Configure instance;

    private Configure() {
        try{
            CONFIG_PATH = searchConfig();
            readConfigureFile(CONFIG_PATH);
        }catch (Exception e) {
            e.printStackTrace();
        }
        
    } 

    public static Configure getInstance(){ 
        if(instance == null) { // 1번 : 쓰레드가 동시 접근시 문제 
            instance = new Configure(); // 2번 : 쓰레드가 동시 접근시 인스턴스 여러번 생성 
        }
        return instance; 
    }

 
    // configure 파일 검색 함수
    public static String searchConfig()  throws FileNotFoundException{
        return System.getenv("DT_CONF_HOME") + "/configure";
    }
    // configure 파일을 읽어서 json 문자열로 변환후 confugureSet 함수 호출
    public static void readConfigureFile(String confFileFullName) throws IOException, FileNotFoundException{

        File jsonFile = new File(confFileFullName);
        FileInputStream jsonStream = new FileInputStream(jsonFile);

        BufferedReader rd = new BufferedReader(new InputStreamReader(jsonStream));
        String line = "";

        StringBuffer response = new StringBuffer();
        while((line = rd.readLine()) != null) {
            response.append(line);
            response.append("\n");
        }
        rd.close();

        String jsonString = response.toString();

        HashMap<String, String> configure_json = jsonStrToMap(jsonString);
        configureSet(configure_json);
    }
    // json 문자열을 Map으로 변환
    public static HashMap<String, String> jsonStrToMap(String jsonStr) {
        Gson gson = new Gson();
        HashMap<String, Object> map = new HashMap<String, Object>();
        map = (HashMap<String, Object>) gson.fromJson(jsonStr, map.getClass());

        HashMap<String, String> return_map = new HashMap<String, String>();
        for (HashMap.Entry<String, Object> entry : map.entrySet()) {
            if(entry.getValue() instanceof String) {
                return_map.put(entry.getKey(), (String) entry.getValue());
            }
            else if(entry.getValue() instanceof Double){
                return_map.put(entry.getKey(), Integer.toString( ((Double)entry.getValue()).intValue() ));
            }
        }
        return return_map;
    }
    // configure 파일 내용을 각 변수로 적용
    public static void configureSet(HashMap<String, String> configure_json) {

      QUEUE_NAME = configure_json.get("queue_name");
      QUEUE_TOPIC = configure_json.get("queue_topic");

      ZOOKEEPER_PORT = configure_json.get("zookeeper_port");

      ZOOKEEPER_QUORUM = configure_json.get("kafkaservers_ip");
      RULE_PATH = configure_json.get("rule_dir");

      MKDRL_PATH = configure_json.get("mkdrl_path");
      //if (MKDRL_PATH.contains(".py")){
       // MKDRL_PATH = "python " + MKDRL_PATH;
      //}
  
    }
}
