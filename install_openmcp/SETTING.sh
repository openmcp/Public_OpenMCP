cp -r master cp_master
cp -r member cp_member

echo -n "OpenmMCP Analytic Engine GRPC Server IP -> "
read OAE_GRPC_IP

echo -n "OpenmMCP Analytic Engine GRPC Server Port -> "
read OAE_GRPC_PORT

sed -i 's|REPLACE_GRPCIP|'"$OAE_GRPC_IP"'|g' master_cp/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'"$OAE_GRPC_IP"'|g' master_cp/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCIP|'"$OAE_GRPC_IP"'|g' master_cp/openmcp-loadbalancing-controller/operator.yaml

sed -i 's|REPLACE_GRPCPORT|'"$OAE_GRPC_PORT"'|g' master_cp/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$OAE_GRPC_PORT"'|g' master_cp/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$OAE_GRPC_PORT"'|g' master_cp/openmcp-loadbalancing-controller/operator.yaml


echo -n "OpenmMCP Metric Collector GRPC Server IP -> "
read OME_GRPC_IP

echo -n "OpenmMCP Metric Collector GRPC Server Port -> "
read OME_GRPC_PORT

sed -i 's|REPLACE_GRPCIP|'"$OME_GRPC_IP"'|g' member_cp/metric-collector/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$OME_GRPC_PORT"'|g' member_cp/metric-collector/operator.yaml

echo -n "InfluxDB Server IP -> "
read INFLUXDB_IP

echo -n "InfluxDB Server Port -> "
read INFLUXDB_PORT

echo -n "InfluxDB User Name -> "
read INFLUXDB_USERNAME

echo -n "InfluxDB User Password -> "
read INFLUXDB_USERPWD

sed -i 's|REPLACE_INFLUXDBIP|'"$INFLUXDB_IP"'|g' master_cp/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'"$INFLUXDB_IP"'|g' master_cp/openmcp-metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'"$INFLUXDB_PORT"'|g' master_cp/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'"$INFLUXDB_PORT"'|g' master_cp/openmcp-metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERNAME|'"$INFLUXDB_USERNAME"'|g' master_cp/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERNAME|'"$INFLUXDB_USERNAME"'|g' master_cp/openmcp-metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERPWD|'"$INFLUXDB_USERPWD"'|g' master_cp/openmcp-analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERPWD|'"$INFLUXDB_USERPWD"'|g' master_cp/openmcp-metric-collector/operator.yaml

echo -n "NFS Server IP -> "
read NFS_IP
sed -i 's|REPLACE_NFSIP|'"$NFS_IP"'|g' master_cp/influxdb/pv.yaml

echo -n "PowerDNS Server IP -> "
read PDNS_IP

echo -n "PowerDNS Server Port -> "
read PDNS_PORT

echo -n "PowerDNS Server API Key -> "
read PDNS_API_KEY

sed -i 's|REPLACE_PDNSIP|'"$PDNS_IP"'|g' master_cp/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSPORT|'"$PDNS_PORT"'|g' master_cp/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'"$PDNS_API_KEY"'|g' master_cp/openmcp-dns-controller/operator.yaml

echo -n "OpenMCP MetalLB Address IP Range (FROM) -> "
read ADDRESS_FROM

echo -n "OpenMCP MetalLB Address IP Range (TO) -> "
read ADDRESS_TO

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master_cp/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master_cp/metallb/configmap.yaml
