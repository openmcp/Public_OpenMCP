cur=`pwd`

cd ../openmcp-analytic-engine
./2.create.sh
cd $cur

cd ../openmcp-apiserver
./2.create.sh
cd $cur

cd ../openmcp-cluster-manager
./2.create.sh
cd $cur

cd ../openmcp-dns-controller
./2.create.sh
cd $cur

cd ../openmcp-loadbalancing-controller
./2.create.sh
cd $cur

cd ../openmcp-metric-collector/master
./2.create.sh
cd $cur

cd ../openmcp-policy-engine
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-configmap-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-daemonset-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-deployment-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-has-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-ingress-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-job-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-namespace-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pv-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pvc-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-secret-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-service-controller
./2.create.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-statefulset-controller
./2.create.sh
cd $cur


cd ../openmcp-scheduler
./2.create.sh
cd $cur

cd ../openmcp-sync-controller
./2.create.sh
cd $cur



