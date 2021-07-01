package com.sample;
//===========json convert============
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;

//===============Drools===============
import org.kie.api.KieBase;
import org.kie.api.KieBaseConfiguration;
import org.kie.api.KieServices;
import org.kie.api.builder.KieBuilder;
import org.kie.api.builder.KieFileSystem;
import org.kie.api.builder.model.KieBaseModel;
import org.kie.api.builder.model.KieModuleModel;
import org.kie.api.builder.model.KieSessionModel;
import org.kie.api.conf.EqualityBehaviorOption;
import org.kie.api.conf.EventProcessingOption;
import org.kie.api.io.Resource;
import org.kie.api.io.ResourceType;
import org.kie.api.runtime.KieContainer;
import org.kie.api.runtime.KieSession;
import org.kie.api.runtime.conf.ClockTypeOption;
import org.kie.api.runtime.rule.FactHandle;

//===========queue===============
import java.util.concurrent.LinkedBlockingQueue;

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

public class RuleEngine implements  Runnable{

    LinkedBlockingQueue<String> queue;
    ArrayList<FactHandle> arraylistFacthandle;
    HashMap<String, String> prev_version_map;
    KieSession kieSession;
    HBaseData hbdata;

    Configure c;

    public RuleEngine(LinkedBlockingQueue<String> queue) throws FileNotFoundException, IOException {
        c = Configure.getInstance();
        this.queue = queue;
        this.arraylistFacthandle = new ArrayList<FactHandle>();
        this.prev_version_map = new HashMap<String, String>();
        this.hbdata = new HBaseData();

    }
    
    public void run() {
        try {

            while(true){
                String message = this.queue.take();

                HashMap<String, String> dataMap = jsonStrToMap(message);
                Sensor object_sensor = new Sensor();
                object_sensor.set_dataMap(dataMap);
                //version_check();
                String CraneFullName = object_sensor.dataMap.get("CraneFullName");
                //String recv_version = object_sensor.dataMap.get("Version");
                String time = object_sensor.dataMap.get("Time");
                
                String current_version =  hbdata.getVersion(CraneFullName);

                String prev_version = this.prev_version_map.get(CraneFullName);
                // System.out.println("******************************");
                // System.out.println("recv_version : "+ recv_version);
                // System.out.println("prev_version : "+ prev_version);
                // System.out.println("current_version : "+ current_version);

                if (prev_version == null || !prev_version.equals(current_version)){
                    Runtime rt = Runtime.getRuntime();
                    //String[] options = new String[]{CraneFullName};
                    String[] command;
                    if (c.MKDRL_PATH.contains(".py")){
                        command = new String[] {"python", c.MKDRL_PATH, CraneFullName};
                    }
    		          else {
                        command = new String[] {c.MKDRL_PATH, CraneFullName};
                    }
                    
                    Process pc = null;
                    try {
                          //pc = rt.exec(c.MKDRL_PATH, options);
                          pc = rt.exec(command);
                          pc.waitFor();
                          
                        } catch (IOException e) {
                          e.printStackTrace();
                        } catch (InterruptedException e){
                          e.printStackTrace();
                        }
                        finally {
                          pc.destroy();
                          prev_version_map.put(CraneFullName, current_version);
                          System.out.println(CraneFullName + " / " + time + " : MakeDRL. (prev_version != current_version)");
                        }
                        //System.out.println("Create KieSession");    
                        kieSession = createKieSession(c.RULE_PATH + "/"+CraneFullName+".drl");
                        arraylistFacthandle.clear();

                        hbdata.sensor_error_load(CraneFullName);
                }
                // if (!recv_version.equals(current_version)){
                //     System.out.println(CraneFullName + " / " + time + " : Ignored. (recv_version != current_version)");
                //     // System.out.println("******************************");
                //     // System.out.println("Version Inconsistency !!");
                //     // System.out.println("Recv Version : " + recv_version + ", Current Version : " + current_version);
                //     // System.out.println("Message ignore");
                    
                //     continue;
                // }
                //System.out.println("******************************");
                object_sensor.set_initResultMap2();
                FactHandle fh = kieSession.insert(object_sensor);

                arraylistFacthandle.add(fh);  
                //System.out.println(kieSession.getFactCount());
                 
                for (int i =0; i< kieSession.getFactCount() - 10; i++){
                    kieSession.retract(arraylistFacthandle.get(0));
                    arraylistFacthandle.remove(0);
                }
                try{
                  kieSession.fireAllRules();
                }
                catch (Exception e){
                  System.out.println("Data Format Error. Ignored");
                  continue;
                }
                hbdata.saveHBaseResult(object_sensor);
                System.out.println(CraneFullName + " / " + time + " : Data Real-time Analysis and Save Complete");
                // System.out.println(time + " : Data Real-time Analysis and Save Complete");
            }

        }
        catch(Exception e) {
            e.printStackTrace();
        }
      
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

    // CEP Engine KieSession 생성
    public KieSession createKieSession(String rule_file) {
      KieServices ks = KieServices.Factory.get();
      // Create a module model
      KieModuleModel kieModuleModel = ks.newKieModuleModel();
      // Base Model from the module model
      KieBaseModel kieBaseModel = kieModuleModel.newKieBaseModel( "KBase" ).setDefault( true ).setEqualsBehavior( EqualityBehaviorOption.EQUALITY).setEventProcessingMode( EventProcessingOption.STREAM );
      //Create session model for the Base Model
      KieSessionModel ksessionModel = kieBaseModel.newKieSessionModel( "KSession" ).setDefault( true ).setType( KieSessionModel.KieSessionType.STATEFUL ).setClockType( ClockTypeOption.get("realtime") );
      // Create File System services
      KieFileSystem kFileSystem = ks.newKieFileSystem();
      // File file = new File("src/resources/rules/Sample.drl");
      File file = new File(rule_file);
      Resource resource = ks.getResources().newFileSystemResource(file).setResourceType(ResourceType.DRL);
      kFileSystem.write( resource );
      KieBuilder kbuilder = ks.newKieBuilder( kFileSystem );
      // kieModule is automatically deployed to KieRepository if successfully built.
      kbuilder.buildAll();
      if (kbuilder.getResults().hasMessages(org.kie.api.builder.Message.Level.ERROR)) {
          throw new RuntimeException("Build time Errors: " + kbuilder.getResults().toString());
      }     
      KieContainer kContainer = ks.newKieContainer(ks.getRepository().getDefaultReleaseId());
      KieBaseConfiguration config = ks.newKieBaseConfiguration();
      config.setOption(EventProcessingOption.STREAM);

      KieBase kieBase = kContainer.newKieBase( config );
      return kieBase.newKieSession();
    }

}
