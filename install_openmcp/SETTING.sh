echo -n "GRPC Server IP -> "
read GRPC_IP

echo -n "GRPC Server Port -> "
read GRPC_PORT

sed -i 's|REPLACE_GRPCIP|'"$GRPC_IP"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'"$GRPC_IP"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCIP|'"$GRPC_IP"'|g' master/loadbalancing-controller/operator.yaml
sed -i 's|REPLACE_GRPCIP|'"$GRPC_IP"'|g' member/metric-collector/operator.yaml

sed -i 's|REPLACE_GRPCPORT|'"$GRPC_PORT"'|g' master/openmcp-has-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$GRPC_PORT"'|g' master/openmcp-scheduler/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$GRPC_PORT"'|g' master/loadbalancing-controller/operator.yaml
sed -i 's|REPLACE_GRPCPORT|'"$GRPC_PORT"'|g' member/metric-collector/operator.yaml

echo -n "InfluxDB Server IP -> "
read INFLUXDB_IP

echo -n "InfluxDB Server Port -> "
read INFLUXDB_PORT

echo -n "InfluxDB User Name -> "
read INFLUXDB_USERNAME

echo -n "InfluxDB User Password -> "
read INFLUXDB_USERPWD

sed -i 's|REPLACE_INFLUXDBIP|'"$INFLUXDB_IP"'|g' master/analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBIP|'"$INFLUXDB_IP"'|g' master/metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBPORT|'"$INFLUXDB_PORT"'|g' master/analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBPORT|'"$INFLUXDB_PORT"'|g' master/metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERNAME|'"$INFLUXDB_USERNAME"'|g' master/analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERNAME|'"$INFLUXDB_USERNAME"'|g' master/metric-collector/operator.yaml

sed -i 's|REPLACE_INFLUXDBUSERPWD|'"$INFLUXDB_USERPWD"'|g' master/analytic-engine/operator.yaml
sed -i 's|REPLACE_INFLUXDBUSERPWD|'"$INFLUXDB_USERPWD"'|g' master/metric-collector/operator.yaml

echo -n "NFS Server IP -> "
read NFS_IP
sed -i 's|REPLACE_NFSIP|'"$NFS_IP"'|g' master/influxdb/pv.yaml

echo -n "PowerDNS Server IP -> "
read PDNS_IP

echo -n "PowerDNS Server Port -> "
read PDNS_PORT

echo -n "PowerDNS Server API Key -> "
read PDNS_API_KEY

sed -i 's|REPLACE_PDNSIP|'"$PDNS_IP"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSPORT|'"$PDNS_PORT"'|g' master/openmcp-dns-controller/operator.yaml
sed -i 's|REPLACE_PDNSAPIKEY|'"$PDNS_API_KEY"'|g' master/openmcp-dns-controller/operator.yaml

echo -n "OpenMCP MetalLB Address IP Range (FROM) -> "
read ADDRESS_FROM

echo -n "OpenMCP MetalLB Address IP Range (TO) -> "
read ADDRESS_TO

sed -i 's|REPLACE_ADDRESS_FROM|'"$ADDRESS_FROM"'|g' master/metallb/configmap.yaml
sed -i 's|REPLACE_ADDRESS_TO|'"$ADDRESS_TO"'|g' master/metallb/configmap.yaml
