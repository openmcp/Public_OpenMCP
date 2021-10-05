USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

#openmcpnamespace 생성 시
#URL="apis/openmcp.k8s.io/v1alpha1/namespaces/default/openmcpnamespaces"

URL="api/v1/namespaces"
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

curl -X POST --cacert server.crt -H 'Content-Type: application/yaml' -H "Authorization: Bearer $TOKEN" --data '
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
' https://$IP:$PORT/$URL?clustername=$CONTEXT

#curl -X POST --cacert server.crt -H 'Content-Type: application/yaml' -H "Authorization: Bearer $TOKEN" --data '
#apiVersion: openmcp.k8s.io/v1alpha1
#kind: OpenMCPNamespace
#metadata:
#  name: test-namespace
#' https://$IP:$PORT/$URL?clustername=$CONTEXT

rm server.crt
