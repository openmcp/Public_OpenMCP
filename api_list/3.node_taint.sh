USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

NODE="openmcp-master"

URL="api/v1/nodes/$NODE"
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

#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data '[{"op": "replace", "path": "/metadata/labels/node-role.kubernetes.io~0master", "value": NONE}]' https://$IP:$PORT/$URL?clustername=$CONTEXT

curl -X PATCH --cacert server.crt -H "Content-Type: application/strategic-merge-patch+json" -H "Authorization: Bearer $TOKEN" \
--data '{"metadata":{"labels":{"node-role.kubernetes.io/master":""}}}' https://$IP:$PORT/$URL?clustername=$CONTEXT

rm server.crt