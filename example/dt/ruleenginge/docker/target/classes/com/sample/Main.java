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


public class Main {
    public static void main(String[] args) {
        Configure c = Configure.getInstance();
        AMQP_Receiver amqp = new AMQP_Receiver();
    }
}

   // public static void make_drl(String drl_path) throws InterruptedException {
    //         Runtime rt = Runtime.getRuntime();
    //         Process pc = null;

    //         try {
    //           pc = rt.exec(MKDRL_PATH);
    //           System.out.println("Mkdrl Excute");
    //         } catch (IOException e) {
    //           e.printStackTrace();
    //         } finally {
    //           pc.waitFor();
    //           pc.destroy();
    //         }

    //         kieSession = createKieSession(drl_path);
    //         // System.out.println("drl upload");

    // }

    // // configure 파일 검색 함수
    // public static String searchConfig() throws FileNotFoundException {
    //   String configFilePath = "";
    //   String cwd = System.getProperty("user.dir");
    //   File f = new File(cwd);
    //   File[] file_list = f.listFiles();

    //   while(true) {
    //     cwd = f.getParent();
    //     if (cwd == null) {
    //       return "-1";
    //     }

    //     f = new File(cwd);
    //     file_list = f.listFiles();
    //     for (File filename : file_list) {
    //       if(filename.getName().equals("conf")) {
    //         if(Files.exists(Paths.get(cwd + "/conf/configure"))) {
    //           configFilePath = cwd + "/conf/configure";
    //           return configFilePath;
    //         }
    //       }
    //     }
    //   }
    // }
     

    
 
    

    
