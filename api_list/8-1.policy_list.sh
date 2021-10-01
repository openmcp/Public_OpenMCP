USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

URL="apis/openmcp.k8s.io/v1alpha1/openmcppolicys"
CONTEXT="openmcp"

echo -n | openssl s_client -connect $IP:$PORT | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > server.crt

TOKEN_JSON=`curl -XPOST \
        --cacert server.crt \
        --insecure \
        -H "Content-type: application/json" \
        --data "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
        https://$IP:$PORT/token`

TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

curl -X GET --cacert server.crt -H "Authorization: Bearer $TOKEN" https://$IP:$PORT/$URL?clustername=$CONTEXT
rm server.crt