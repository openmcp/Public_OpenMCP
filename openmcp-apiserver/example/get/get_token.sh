USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
URLPATH="token"

curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/$URLPATH?username=$USERNAME&password=$PASSWORD"

