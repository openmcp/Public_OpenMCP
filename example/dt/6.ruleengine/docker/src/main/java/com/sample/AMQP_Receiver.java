package com.sample;

//===========json convert============
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;

// ==========rabbitmq import===============
import com.rabbitmq.client.AMQP;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.Consumer;
import com.rabbitmq.client.DefaultConsumer;
import com.rabbitmq.client.Envelope;

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

public class AMQP_Receiver{
    HashMap<String, String> dataMap;
    HashMap<String, RuleEngine> threadMap;
    HashMap<String ,LinkedBlockingQueue<String>> queueMap;
    public static Configure c;

    public AMQP_Receiver(){
        c = Configure.getInstance();
        if (c.QUEUE_NAME.equals("data")) {
            System.out.println("Start Rule Engine");
        }
        threadMap = new HashMap<String, RuleEngine>();
        queueMap = new HashMap<String, LinkedBlockingQueue<String>>();
 
        waitMessage();
    }

    
    // RabbitMQ 설정 함수
    public static Channel createQueueConfigure() throws IOException, TimeoutException {
      ConnectionFactory factory = new ConnectionFactory();

      factory.setHost(c.RABBITMQ_SERVICE);
      factory.setUsername("dtuser01");
      factory.setPassword("dtuser01");
      Connection connection = factory.newConnection();
      Channel channel = connection.createChannel();
      return channel;
    }

    public void waitMessage() {
        try{
            Channel channel = createQueueConfigure();
            channel.queueDeclare(c.QUEUE_NAME, false, false, false, null);

            System.out.println(" [*] Waiting for messages. To exit press CTRL+C");

            Consumer consumer = new DefaultConsumer(channel) {
                @Override
                public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body) throws IOException {
                    //System.out.println("Msg Recv !! push queue");
                    String message = new String(body, "UTF-8");
                    dataMap = jsonStrToMap(message);

                    String CraneFullName = dataMap.get("CraneFullName");
                    LinkedBlockingQueue<String> queue;
                    if (threadMap.containsKey(CraneFullName)) { // key가 포함되어있으면
                        //RuleEngine re = threadMap.get(CraneFullName);
                        queue = queueMap.get(CraneFullName);


                    }
                    else{
                        queue = new LinkedBlockingQueue<String>();
                                            
                        RuleEngine re = new RuleEngine(queue);
                        threadMap.put(CraneFullName, re);
                        queueMap.put(CraneFullName, queue);

                        Thread t = new Thread(re);
                        t.start();
                        System.out.println("Thread '"+ CraneFullName +"' Start !!");

                    }
                    try{
                        queue.put(message);
                    } catch(Exception e){
                        e.printStackTrace();
                    }
                }
        
            };
        channel.basicConsume(c.QUEUE_TOPIC, true, consumer);
        }catch (Exception e) {
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

}
