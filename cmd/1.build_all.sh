cur=`pwd`

#kubectl apply -f ../crds

cd ../openmcp-analytic-engine
./1.build.sh
cd $cur

cd ../openmcp-apiserver
./1.build.sh
cd $cur

cd ../openmcp-cluster-manager
./1.build.sh
cd $cur

cd ../openmcp-dns-controller
./1.build.sh
cd $cur

cd ../openmcp-loadbalancing-controller
./1.build.sh
cd $cur

cd ../openmcp-metric-collector/master
./1.build.sh
cd $cur

cd ../openmcp-policy-engine
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-configmap-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-daemonset-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-deployment-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-has-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-ingress-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-job-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-namespace-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pv-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-pvc-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-secret-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-service-controller
./1.build.sh
cd $cur

cd ../openmcp-resource-controller/openmcp-statefulset-controller
./1.build.sh
cd $cur


cd ../openmcp-scheduler
./1.build.sh
cd $cur

cd ../openmcp-sync-controller
./1.build.sh
cd $cur



