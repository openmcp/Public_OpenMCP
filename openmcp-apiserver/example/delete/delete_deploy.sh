USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
NS="default"
POD="example-nginx-deploy"
URL="apis/apps/v1/namespaces/$NS/deployments/$POD"
CLUSTER="openmcp"

TOKEN_JSON=`curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/token?username=$USERNAME&password=$PASSWORD"`
TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

curl -X DELETE -H "Authorization: Bearer $TOKEN" $IP:$PORT/$URL?clustername=$CLUSTER

