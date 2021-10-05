USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

NS="default"
PODNAME="test-deploy"
CONTAINERNAME="nginx"
CPU="200m"
MEM="20Mi"

URL="apis/apps/v1/namespaces/$NS/deployments/$PODNAME"
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


curl -X PATCH --cacert server.crt -H "Content-Type: application/strategic-merge-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"$CONTAINERNAME\", \"resources\":{\"requests\":{\"memory\":\"$MEM\", \"cpu\":\"$CPU\"}}}]}}}}" https://$IP:$PORT/$URL?clustername=$CONTEXT


rm server.crt