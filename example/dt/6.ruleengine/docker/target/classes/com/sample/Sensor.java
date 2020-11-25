package com.sample;

import java.util.HashMap;
import java.util.List;

public class Sensor {

    HashMap<String, String> dataMap = new HashMap();
    HashMap<String, String> resultMap = new HashMap();

    public void set_dataMap(HashMap<String, String> dataMap) {
        this.dataMap = dataMap;
    }

    public HashMap<String, String> get_dataMap() {
        return this.dataMap;
    }

    public void set_initResultMap(List<String> sensorList) {
        for(String sensor : sensorList ) {
            this.set_resultMap(sensor, "O");
        }
    }
    public void set_initResultMap2(){
        for(HashMap.Entry<String, String> entry : this.dataMap.entrySet()) {
            String sensor_name = entry.getKey();
            if (!sensor_name.equals("CraneFullName") && !sensor_name.equals("Version") && !sensor_name.equals("Time")){
                this.set_resultMap(entry.getKey(), "O");
            }
            
        }
    }

    public void set_resultMap(String key, String value) {
        this.resultMap.put(key, value);
    }

    public HashMap<String, String> get_resultMap() {
        return this.resultMap;
    }

    public int error_value = 0;

    public String get_ruleErrorCode() {
        //return ("0x" + this.error_value.toHexString);
        return ("0x" + Integer.toHexString(this.error_value));
    }

    public void set_ruleErrorCode(String value) {
        //this.error_value = this.error_value | Math.pow(2, value.toInt - 1).toInt;
        this.error_value = this.error_value | (int)(Math.pow(2, Integer.parseInt(value) - 1));
    }

}
