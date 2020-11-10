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
    public static String RABBITMQ_SERVICE = "";
    public static String MKDRL_PATH = "";
    public static String RULE_PATH ="";
    public static String INFLUXDB = "";
    public static String INFLUX_PORT = "";
    public static String INFLUX_USER = "";
    public static String INFLUX_PW = "";




    private static Configure instance;

    private Configure() {
        RABBITMQ_SERVICE = System.getenv("RABBITMQ_SERVICE");
        QUEUE_NAME = System.getenv("QUEUE_NAME");
        QUEUE_TOPIC = System.getenv("QUEUE_TOPIC");
        MKDRL_PATH = System.getenv("MKDRL_PATH");
        RULE_PATH = System.getenv("RULE_DIR");
        INFLUXDB = System.getenv("INFLUXDB");
        INFLUX_PORT = System.getenv("INFLUX_PORT");
        INFLUX_USER = System.getenv("INFLUX_USER");
        INFLUX_PW = System.getenv("INFLUX_PW");


	    System.out.println("Config Complete");
        System.out.println(QUEUE_NAME);
        
    } 

    public static Configure getInstance(){ 
        if(instance == null) { // 1번 : 쓰레드가 동시 접근시 문제 
            instance = new Configure(); // 2번 : 쓰레드가 동시 접근시 인스턴스 여러번 생성 
        }
        return instance; 
    }


}
