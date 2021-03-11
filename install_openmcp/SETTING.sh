if [ -d "master" ]; then
  # Control will enter here if $DIRECTORY exists.
  rm -r master
fi

if [ -d "member" ]; then
  # Control will enter here if $DIRECTORY exists.
  rm -r member
fi

cp -r master.back master
cp -r member.back member

#MYIP=`ip route get 8.8.8.8 | head -1 | cut -d' ' -f8`

echo -n "OpenMCP Install Type [debug/learning]-> "
read OMCP_INSTALL_TYPE

echo -n "OpenMCP Server IP -> "
read OMCP_IP

echo -n "Docker Secret Name for Authentication -> "
read DOCKER_SECRET_NAME

echo -n "Docker Registry IP:PORT -> "
read DOCKER_REGISTRY_IP

echo -n "Docker imagePullPolicy [Always/IfNotPresent] -> "
read DOCKER_IMAGE_PULL_POLICY

#echo -n "OpenMCP Analytic Engine GRPC Server IP -> "
#read OAE_GRPC_IP

echo -n "OpenMCP API Server Port -> "
read OAS_PORT

echo -n "OpenMCP Cluster Manager Port -> "
read OCM_PORT

echo -n "OpenMCP Analytic Engine GRPC Server Port -> "
read OAE_GRPC_PORT

echo -n "OpenMCP Metric Collector GRPC Server IP (Public) -> "
read OME_GRPC_PUBLIC_IP

echo -n "OpenMCP Metric Collector GRPC Server Port (Public) -> "
read OME_GRPC_PUBLIC_PORT

#echo -n "InfluxDB Server IP -> "
#read INFLUXDB_IP

echo -n "InfluxDB Server Port -> "
read INFLUXDB_PORT

echo -n "InfluxDB User Name -> "
read INFLUXDB_USERNAME

echo -n "InfluxDB User Password -> "
read INFLUXDB_USERPWD

echo -n "NFS & PowerDNS Server IP -> "
read NFS_PDNS_IP

echo -n "PowerDNS Server IP (Public) -> "
read PDNS_PUBLIC_IP

echo -n "PowerDNS Server Port (public) -> "
read PDNS_PUBLIC_PORT

echo -n "PowerDNS Server API Key -> "
read PDNS_API_KEY

echo -n "OpenMCP MetalLB Address IP Range (FROM) -> "
read ADDRESS_FROM

echo -n "OpenMCP MetalLB Address IP Range (TO) -> "
read ADDRESS_TO

if [ $OMCP_INSTALL_TYPE == "learning" ]; then
  rm master/openmcp-cluster-manager/operator.yaml
  rm master/influxdb/deployment.yaml
  mv master/openmcp-cluster-manager/operator-learningmcp.yaml master/openmcp-cluster-manager/operator.yaml
  mv master/influxdb/deployment-learningmcp.yaml master/influxdb/deployment.yaml
else
  rm master/openmcp-cluster-manager/operator-learningmcp.yaml
  rm master/influxdb/deployment-learningmcp.yaml
fi

sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-loadbalancing-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-sync-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-configmap-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-ingress-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-secret-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-deployment-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-service-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-policy-engine/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-loadbalancing-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-sync-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-configmap-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-ingress-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-secret-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-deployment-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-service-controller/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-policy-engine/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' master/openmcp-cluster-manager/operator.yaml
sed -i 's|REPLACE_DOCKERREGISTRYIP|'$DOCKER_REGISTRY_IP'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/influxdb/deployment.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-loadbalancing-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-sync-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-configmap-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-ingress-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-secret-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-deployment-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-service-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-policy-engine/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-cluster-manager/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_OMCPIP|'\"$OMCP_IP\"'|g' master/openmcp-cluster-manager/pv.yaml

sed -i 's|REPLACE_PORT|'$OAS_PORT'|g' master/openmcp-apiserver/service.yaml
sed -i 's|REPLACE_PORT|'$OCM_PORT'|g' master/openmcp-cluster-manager/service.yaml

sed -i 's|REPLACE_GRPCPORT|'$OAE_GRPC_PORT'|g' master/openmcp-analytic-engine/service.yaml

sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_GRPCIP|'\"$OME_GRPC_PUBLIC_IP\"'|g' member/metric-collector/operator/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'$OME_GRPC_PUBLIC_PORT'|g' member/metric-collector/operator/operator.yaml
#sed -i 's|REPLACE_GRPCPORT|'$OME_GRPC_PUBLIC_PORT'|g' member/metric-collector/service.yaml
sed -i 's|REPLACE_GRPCPORT|'$OME_GRPC_PUBLIC_PORT'|g' master/openmcp-metric-collector/service.yaml

sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'$INFLUXDB_PORT'|g' master/influxdb/service.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERNAME|'\"$INFLUXDB_USERNAME\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERNAME|'\"$INFLUXDB_USERNAME\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERNAME|'\"$INFLUXDB_USERNAME\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERPWD|'\"$INFLUXDB_USERPWD\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERPWD|'\"$INFLUXDB_USERPWD\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERPWD|'\"$INFLUXDB_USERPWD\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_NFSIP|'\"$NFS_PDNS_IP\"'|g' master/influxdb/pv.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' master/configmap/coredns/coredns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/coredns/coredns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/kubedns/kube-dns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'\"$NFS_PDNS_IP\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'\"$PDNS_API_KEY\"'|g' master/openmcp-dns-controller/operator.yaml

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master/metallb/configmap.yaml

USERNAME=`whoami`

if [ $OMCP_INSTALL_TYPE == "learning" ]; then
  cp -r member /home/$USERNAME/.init/member
fi