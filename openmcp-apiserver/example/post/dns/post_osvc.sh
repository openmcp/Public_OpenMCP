USERNAME="openmcp"
PASSWORD="keti"
#IP="10.0.3.20"
#PORT="30000"
IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"
NS="default"
URL="apis/openmcp.k8s.io/v1alpha1/namespaces/$NS/openmcpservices"
CLUSTER="openmcp"

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
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPService
metadata:
  name: openmcp-service-example
spec:
  template:
    spec:
      type: LoadBalancer
      ports:
        - name: http
          port: 80
      selector:
        app: nginx
' https://$IP:$PORT/$URL?clustername=$CLUSTER
rm server.crt
