kubectl delete secret influxdb-creds --context cluster3
kubectl delete -f deploy --context cluster3
kubectl delete -f sdr.yaml
