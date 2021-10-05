#!/bin/bash


# REQUIRED_PKG=python-pip
# PKG_OK=$(dpkg-query -W --showformat='${Status}\n' $REQUIRED_PKG|grep "install ok installed")
# if [ "" = "$PKG_OK" ]; then
#   echo "No $REQUIRED_PKG. Setting up $REQUIRED_PKG."
#   sudo apt-get --yes install $REQUIRED_PKG 
# fi

#REQUIRED_PKG=python-pip
#PKG_OK=$(dpkg-query -W --showformat='${Status}\n' $REQUIRED_PKG|grep "install ok installed")
#if [ "" = "$PKG_OK" ]; then
#  echo "No $REQUIRED_PKG. Setting up $REQUIRED_PKG."
#  sudo apt-get --yes install $REQUIRED_PKG 
#fi

REQUIRED_PKG2=nfs-kernel-server
PKG_OK2=$(dpkg-query -W --showformat='${Status}\n' $REQUIRED_PKG2|grep "install ok installed")
if [ "" = "$PKG_OK2" ]; then
  echo "No $REQUIRED_PKG2. Setting up $REQUIRED_PKG2."
  sudo apt-get --yes install $REQUIRED_PKG2 
fi


curl -O https://bootstrap.pypa.io/pip/2.7/get-pip.py
python get-pip.py
rm get-pip.py


PYTHON_REQUIRED_PKG=yq
PKG_OK=$(pip list --disable-pip-version-check | grep $PYTHON_REQUIRED_PKG)
if [ "" = "$PKG_OK" ]; then
  echo "No $PYTHON_REQUIRED_PKG. Setting up $PYTHON_REQUIRED_PKG."
  sudo pip install $PYTHON_REQUIRED_PKG 
fi


CONFFILE=settings.yaml

OMCP_INSTALL_TYPE=`yq -r .default.installType $CONFFILE`

DOCKER_SECRET_NAME=`yq -r .default.docker.imagePullSecretName $CONFFILE`
DOCKER_IMAGE_PULL_POLICY=`yq -r .default.docker.imagePullPolicy $CONFFILE`

OMCP_IP=`yq -r .master.ServerIP.internal $CONFFILE`
OMCP_EXTERNAL_IP=`yq -r .master.ServerIP.external $CONFFILE`

ADDRESS_FROM=`yq -r .master.metalLB.rangeStartIP $CONFFILE`
ADDRESS_TO=`yq -r .master.metalLB.rangeEndIP $CONFFILE`

OAS_NODE_PORT=`yq -r .master.Moudules.APIServer.NodePort $CONFFILE`
API_APP_KEY=`yq -r .master.Moudules.APIServer.AppKey $CONFFILE`
API_USER_NAME=`yq -r .master.Moudules.APIServer.UserName $CONFFILE`
API_USER_PW=`yq -r .master.Moudules.APIServer.UserPW $CONFFILE`

OAE_NODE_PORT=`yq -r .master.Moudules.AnalyticEngine.NodePort $CONFFILE`

OME_NODE_PORT=`yq -r .master.Moudules.MetricCollector.NodePort $CONFFILE`
OME_EXTERNAL_PORT=`yq -r .master.Moudules.MetricCollector.externalPort $CONFFILE`

INFLUXDB_NODE_PORT=`yq -r .master.Moudules.InfluxDB.NodePort $CONFFILE`

LB_EXTERNAL_IP=`yq -r .master.Moudules.LoadBalancingController.external $CONFFILE`
LB_NODE_PORT=`yq -r .master.Moudules.LoadBalancingController.NodePort $CONFFILE`

PDNS_IP=`yq -r .externalServer.ServerIP.internal $CONFFILE`
PDNS_PUBLIC_IP=`yq -r .externalServer.ServerIP.external $CONFFILE`
PDNS_PUBLIC_PORT=`yq -r .externalServer.powerDNS.externalPort $CONFFILE`
PDNS_API_KEY=`yq -r .externalServer.powerDNS.apiKey $CONFFILE`




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

# Init Memeber Dir NFS Setting
INIT_MEMBER_DIR=`pwd`/member

NFS_OK=$(grep -r '/root/.kube' /etc/exports)
if [ "" = "$NFS_OK" ]; then
  echo "Not found NFS Setting. Add Export '/root/.kube' in /etc/exports"
  echo "/root/.kube *(rw,no_root_squash,sync,no_subtree_check)" >> /etc/exports
fi

NFS_OK2=$(grep -r $INIT_MEMBER_DIR /etc/exports)
if [ "" = "$NFS_OK2" ]; then
  echo "Not found NFS Setting. Add Export '$INIT_MEMBER_DIR' in /etc/exports"
  echo "$INIT_MEMBER_DIR *(rw,no_root_squash,sync,no_subtree_check)" >> /etc/exports
fi

NFS_OK3=$(grep -r '/home/nfs' /etc/exports)
if [ "" = "$NFS_OK3" ]; then
  echo "Not found NFS Setting. Add Export '/home/nfs' in /etc/exports"
  echo "/home/nfs *(rw,no_root_squash,sync,no_subtree_check)" >> /etc/exports
fi
exportfs -a

# Init /etc/resolv.conf
DNS_OK=$(grep -r "nameserver ${PDNS_IP}" /etc/resolv.conf)
if [ "" = "$DNS_OK" ]; then
  echo "Not found External DNS Server. Add 'nameserver ${PDNS_IP}' in /etc/resolv.conf"
  sed -i "1s/^/nameserver ${PDNS_IP}\n /" /etc/resolv.conf
fi


echo "Replace Setting Variable"
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
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-namespace-controller/operator.yaml
sed -i 's|REPLACE_DOCKERSECRETNAME|'\"$DOCKER_SECRET_NAME\"'|g' master/openmcp-job-controller/operator.yaml
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
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-namespace-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' master/openmcp-job-controller/operator.yaml
sed -i 's|REPLACE_DOCKERIMAGEPULLPOLICY|'$DOCKER_IMAGE_PULL_POLICY'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$OMCP_IP\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_INIT_MEMBER_DIR|'\"$INIT_MEMBER_DIR\"'|g' master/openmcp-cluster-manager/pv.yaml 
sed -i 's|REPLACE_OMCPIP|'\"$OMCP_IP\"'|g' master/openmcp-cluster-manager/pv.yaml
sed -i 's|REPLACE_OMCPIP|'\"$OMCP_IP\"'|g' master/openmcp-apiserver/pv.yaml

sed -i 's|REPLACE_PORT|'$OAS_NODE_PORT'|g' master/openmcp-apiserver/service.yaml

sed -i 's|REPLACE_GRPCPORT|'$OAE_NODE_PORT'|g' master/openmcp-analytic-engine/service.yaml

sed -i 's|REPLACE_GRPCPORT|'\"$OAE_NODE_PORT\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_NODE_PORT\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_NODE_PORT\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_GRPCIP|'\"$OMCP_EXTERNAL_IP\"'|g' member/metric-collector/operator/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OME_EXTERNAL_PORT\"'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_GRPCPORT|'$OME_NODE_PORT'|g' master/openmcp-metric-collector/service.yaml

sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$OMCP_IP\"'|g' master/openmcp-cluster-manager/operator.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'$INFLUXDB_NODE_PORT'|g' master/influxdb/service.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_NODE_PORT\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_NODE_PORT\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_NODE_PORT\"'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'\"$INFLUXDB_NODE_PORT\"'|g' master/openmcp-cluster-manager/operator.yaml

sed -i 's|REPLACE_API_KEY|'\"$API_APP_KEY\"'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_API_USER_NAME|'\"$API_USER_NAME\"'|g' master/openmcp-apiserver/operator.yaml
sed -i 's|REPLACE_API_USER_PW|'\"$API_USER_PW\"'|g' master/openmcp-apiserver/operator.yaml

sed -i 's|REPLACE_EXTERNAL_IP|'\"$LB_EXTERNAL_IP\"'|g' master/openmcp-ingress-controller/operator.yaml
sed -i 's|REPLACE_PORT|'$LB_NODE_PORT'|g' master/openmcp-loadbalancing-controller/service.yaml

sed -i 's|REPLACE_NFSIP|'\"$OMCP_IP\"'|g' master/influxdb/pv.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' master/configmap/coredns/coredns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/coredns/coredns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' master/configmap/kubedns/kube-dns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'$PDNS_PUBLIC_IP':'$PDNS_PUBLIC_PORT'|g' member/configmap/kubedns/kube-dns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'\"$PDNS_IP\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'\"$PDNS_API_KEY\"'|g' master/openmcp-dns-controller/operator.yaml

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master/metallb/configmap.yaml

echo "Replace Setting Variable Complete"
USERNAME=`whoami`

if [ $OMCP_INSTALL_TYPE == "learning" ]; then
  echo "Copy 'member' directory"
  rm -rf /home/$USERNAME/.init/member
  cp -r member /home/$USERNAME/.init/member
fi

chmod 755 master/*.sh
chmod 755 member/istio/*.sh


echo "Complete Make Dir(master/member)" 
