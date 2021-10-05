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


# 마스터에 포드 할당 허용
curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data '[{"op": "remove", "path": "/spec/taints/0"}]' https://$IP:$PORT/$URL?clustername=$CONTEXT

# 마스터에 포드 할당 불가
curl -X PATCH --cacert server.crt -H "Content-Type: application/strategic-merge-patch+json" -H "Authorization: Bearer $TOKEN" \
--data '{"spec":{"taints":[{"effect":"NoSchedule", "key":"node-role.kubernetes.io/master"}]}}' https://$IP:$PORT/$URL?clustername=$CONTEXT

rm server.crt