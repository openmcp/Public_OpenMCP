kubectl create secret  generic influxdb-creds --context cluster3 \
--from-literal=INFLUXDB_DATABASE=mydb \
--from-literal=INFLUXDB_USERNAME=root \
--from-literal=INFLUXDB_PASSWORD=root \
--from-literal=INFLUXDB_HOST=influxd

