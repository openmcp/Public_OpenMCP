USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

UNJOINCLUSTER="cluster2"

URL="apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/$UNJOINCLUSTER"
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
--data '[{"op": "replace", "path": "/spec/metalLBRange/addressFrom", "value": ""},{"op": "replace", "path": "/spec/metalLBRange/addressTo", "value": ""},{"op": "replace", "path": "/spec/joinStatus", "value": "UNJOIN"}]' https://$IP:$PORT/$URL?clustername=$CONTEXT


rm server.crt