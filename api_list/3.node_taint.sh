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


# 노드에 포드 할당 허용
# 해당 노드의 taint가 하나인 경우 INDEX = 0 이지만, 여러 개인 경우 "NoSchedule" taint 인덱스 값을 받아오는 과정 필요
# "NoSchedule" taint가 설정되어 있지 않은데 삭제하려고 하는 경우 에러 출력
INDEX=$(kubectl get node $NODE -o json  | jq '.spec.taints | map(.effect == "NoSchedule") | index(true)')
curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "[{\"op\": \"remove\", \"path\": \"/spec/taints/$INDEX\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT



# 노드에 포드 할당 불가
#curl -X PATCH --cacert server.crt -H "Content-Type: application/strategic-merge-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data '{"spec":{"taints":[{"effect":"NoSchedule", "key":"node-role.kubernetes.io/master"}]}}' https://$IP:$PORT/$URL?clustername=$CONTEXT

rm server.crt