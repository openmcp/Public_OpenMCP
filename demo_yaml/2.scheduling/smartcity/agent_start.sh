IP="web.smartcity.openmcp.com"
PORT="8181"
curl --location --request GET "http://$IP:$PORT/restApi/pushConf/AgentMqttt/mqttOffstreet" --header "Content-Type: application/json"
curl --location --request GET "http://$IP:$PORT/restApi/pushConf/AgentWeatherObserved/weatherObserved" --header "Content-Type: application/json"

