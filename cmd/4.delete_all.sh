cur=`pwd`

cd ../openmcp-analytic-engine
./4.delete.sh
cd $cur

cd ../openmcp-apiserver
./4.delete.sh
cd $cur

cd ../openmcp-cluster-manager
./4.delete.sh
cd $cur

cd ../openmcp-dns-controller
./4.delete.sh
cd $cur

cd ../openmcp-loadbalancing-controller
./4.delete.sh
cd $cur

cd ../openmcp-metric-collector/master
./4.delete.sh
cd $cur

cd ../openmcp-policy-engine
./4.delete.sh
cd $cur


cd ../openmcp-resource-controller/openmcp-configmap-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-daemonset-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-deployment-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-has-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-ingress-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-job-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-namespace-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pv-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pvc-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-secret-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-service-controller
./4.delete.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-statefulset-controller
./4.delete.sh
cd $cur


cd ../openmcp-scheduler
./4.delete.sh
cd $cur

cd ../openmcp-sync-controller
./4.delete.sh
cd $cur



