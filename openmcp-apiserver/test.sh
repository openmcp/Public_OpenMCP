USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
URL="api"
CLUSTER="cluster1"

TOKEN_JSON=`curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/token?username=$USERNAME&password=$PASSWORD"`
TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

for ((;;)){
	curl -XGET -H "Authorization: Bearer $TOKEN" $IP:$PORT/$URL?clustername=$CLUSTER
	sleep 1
}
