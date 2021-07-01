#!/bin/bash
pip install yq

CONFFILE=settings.yaml

OMCP_INSTALL_TYPE=`yq -r .default.installType $CONFFILE`

DOCKER_SECRET_NAME=`yq -r .default.docker.imagePullSecretName $CONFFILE`
DOCKER_IMAGE_PULL_POLICY=`yq -r .default.docker.imagePullPolicy $CONFFILE`

OMCP_IP=`yq -r .master.internal.ip $CONFFILE`
OAS_PORT=`yq -r .master.internal.ports.apiServerPort $CONFFILE`
OCM_PORT=`yq -r .master.internal.ports.clusterManagerPort $CONFFILE`
OAE_GRPC_PORT=`yq -r .master.internal.ports.analyticEnginePort $CONFFILE`


OME_GRPC_PUBLIC_IP=`yq -r .master.public.ip $CONFFILE`
OME_GRPC_PUBLIC_PORT=`yq -r .master.public.ports.metricCollectorPort $CONFFILE`

INFLUXDB_PORT=`yq -r .master.internal.ports.influxDBPort $CONFFILE`


PDNS_IP=`yq -r .powerDNS.internal.ip $CONFFILE`
PDNS_PUBLIC_IP=`yq -r .powerDNS.public.ip $CONFFILE`
PDNS_PUBLIC_PORT=`yq -r .powerDNS.public.Ports.pdnsPort $CONFFILE`
PDNS_API_KEY=`yq -r .powerDNS.apiKey $CONFFILE`

ADDRESS_FROM=`yq -r .master.metalLB.rangeStartIP $CONFFILE`
ADDRESS_TO=`yq -r .master.metalLB.rangeEndIP $CONFFILE`


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

if [ $OMCP_INSTALL_TYPE == "learning" ]; then
  rm master/openmcp-cluster-manager/operator.yaml
  rm master/influxdb/deployment.yaml
  mv master/openmcp-cluster-manager/operator-learningmcp.yaml master/openmcp-cluster-manager/operator.yaml
  mv master/influxdb/deployment-learningmcp.yaml master/influxdb/deployment.yaml
else
  rm master/openmcp-cluster-manager/operator-learningmcp.yaml
  rm master/influxdb/deployment-learningmcp.yaml
fi

sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/1.create.sh
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
sed -i 's|REPLACE_GRPCPORT|'\"$OME_GRPC_PUBLIC_PORT\"'|g' member/metric-collector/operator/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'$OME_GRPC_PUBLIC_PORT'|g' master/openmcp-metric-collector/service.yaml

sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'$INFLUXDB_PORT'|g' master/influxdb/service.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_PORT\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_NFSIP|'\"$OMCP_IP\"'|g' master/influxdb/pv.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' master/configmap/coredns/coredns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/coredns/coredns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' master/configmap/kubedns/kube-dns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/kubedns/kube-dns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'\"$PDNS_IP\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'\"$PDNS_API_KEY\"'|g' master/openmcp-dns-controller/operator.yaml

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master/metallb/configmap.yaml

USERNAME=`whoami`

if [ $OMCP_INSTALL_TYPE == "learning" ]; then
  echo "Copy 'member' directory"
  rm -rf /home/$USERNAME/.init/member
  cp -r member /home/$USERNAME/.init/member
fi

mkdir -p /home/nfs/pv/influxdb
