USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

JOINCLUSTER="cluster2"
ADDRESSFROM="10.0.3.211"
ADDRESSTO="10.0.3.220"


URL="apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/$JOINCLUSTER"
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

curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data '[{"op": "replace", "path": "/spec/joinStatus", "value": "JOINING"}]' https://$IP:$PORT/$URL?clustername=$CONTEXT

curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "[{\"op\": \"replace\", \"path\": \"/spec/metalLBRange/addressFrom\", \"value\": \"$ADDRESSFROM\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT

curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "[{\"op\": \"replace\", \"path\": \"/spec/metalLBRange/addressTo\", \"value\": \"$ADDRESSTO\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


rm server.crt