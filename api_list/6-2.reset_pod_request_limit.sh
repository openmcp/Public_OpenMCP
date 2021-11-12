USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

NS="default"
PODNAME="test-deploy"
CONTAINERNAME="nginx"


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

# container가 하나인 경우 INDEX = 0 이지만, 여러 개인 경우 request,limit을 삭제하고자 하는 특정 컨테이너의 인덱스 값을 받아오는 과정 필요
INDEX=$(kubectl get deployment $PODNAME -n $NS -o json  | jq ".spec.template.spec.containers | map(.name == \"$CONTAINERNAME\") | index(true)")

curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "[{\"op\": \"remove\", \"path\": \"/spec/template/spec/containers/$INDEX/resources\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


rm server.crt