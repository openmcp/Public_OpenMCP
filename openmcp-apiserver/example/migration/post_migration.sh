USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="31635"
URL="apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/migrations"
CLUSTER="openmcp"

TOKEN_JSON=`curl -XGET -H "Content-type: application/json" "http://$IP:$PORT/token?username=$USERNAME&password=$PASSWORD"`
TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

curl -X POST -H 'Content-Type: application/yaml' -H "Authorization: Bearer $TOKEN" --data '
apiVersion: openmcp.k8s.io/v1alpha1
kind: Migration
metadata:
  name: migrations
  namespace: openmcp
spec:
  MigrationServiceSource:
  - SourceCluster: cluster1
    TargetCluster: cluster2
    NameSpace: openmcp
    ServiceName: testim11
    MigrationSource:
    - ResourceName: testim-dp
      ResourceType: Deployment
    - ResourceName: testim-sv
      ResourceType: Service
    - ResourceName: testim-pv
      ResourceType: PersistentVolume
    - ResourceName: testim-pvc
      ResourceType: PersistentVolumeClaim
' $IP:$PORT/$URL?clustername=$CLUSTER

