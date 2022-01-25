IP=52.225.33.8
curl --location --request GET "http://$IP:8080/restApi/pushConf/AgentMqttt/mqttOffstreet" --header "Content-Type: application/json"
curl --location --request GET "http://$IP:8080/restApi/pushConf/AgentWeatherObserved/weatherObserved" --header "Content-Type: application/json"
