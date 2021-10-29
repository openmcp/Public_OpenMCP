USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

POD="test-namespace"

#openmcpnamespace 생성 시
#URL="apis/openmcp.k8s.io/v1alpha1/namespaces/default/openmcpnamespaces/$POD"

URL="api/v1/namespaces/$POD"
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

curl -X DELETE --cacert server.crt -H "Authorization: Bearer $TOKEN" https://$IP:$PORT/$URL?clustername=$CONTEXT

rm server.crt
