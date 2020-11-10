kubectl create secret generic influxdb-creds \
--from-literal=INFLUXDB_DATABASE=mydb \
--from-literal=INFLUXDB_USERNAME=root \
--from-literal=INFLUXDB_PASSWORD=root \
--from-literal=INFLUXDB_HOST=influxd

