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

MYIP=`ip route get 8.8.8.8 | head -1 | cut -d' ' -f8`

echo -n "OpenmMCP Analytic Engine GRPC Server Port -> "
read OAE_GRPC_PORT

echo -n "OpenmMCP Metric Collector GRPC Server Port -> "
read OME_GRPC_PORT

echo -n "InfluxDB Server Port -> "
read INFLUXDB_PORT

echo -n "InfluxDB User Name -> "
read INFLUXDB_USERNAME

echo -n "InfluxDB User Password -> "
read INFLUXDB_USERPWD

echo -n "NFS & PowerDNS Server IP -> "
read NFS_PDNS_IP

echo -n "PowerDNS Server Port -> "
read PDNS_PORT

echo -n "PowerDNS Server API Key -> "
read PDNS_API_KEY

echo -n "OpenMCP MetalLB Address IP Range (FROM) -> "
read ADDRESS_FROM

echo -n "OpenMCP MetalLB Address IP Range (TO) -> "
read ADDRESS_TO

sed -i 's|REPLACE_GRPCIP|'\"$MYIP\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$MYIP\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCIP|'\"$MYIP\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OAE_GRPC_PORT\"'|g' master/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_GRPCIP|'\"keti.asuscomm.com\"'|g' member/metric-collector/operator/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'\"$OME_GRPC_PORT\"'|g' member/metric-collector/operator/operator.yaml

sed -i 's|REPLACE_INFLUXDBIP|'\"$MYIP\"'|g' master/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$MYIP\"'|g' master/openmcp-metric-collector/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'\"$MYIP\"'|g' master/openmcp-apiserver/operator.yaml

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

sed -i 's|REPLACE_PDNSIP|'\"$NFS_PDNS_IP\"'|g' master/configmap/coredns/coredns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'\"$NFS_PDNS_IP\"'|g' member/configmap/coredns/coredns-cm.yaml
sed -i 's|REPLACE_PDNSIP|'\"$NFS_PDNS_IP\"'|g' member/configmap/kubedns/kube-dns-cm.yaml

sed -i 's|REPLACE_PDNSIP|'\"$NFS_PDNS_IP\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSPORT|'\"$PDNS_PORT\"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'\"$PDNS_API_KEY\"'|g' master/openmcp-dns-controller/operator.yaml

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master/metallb/configmap.yaml
