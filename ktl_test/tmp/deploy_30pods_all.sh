#!/bin/bash

clusterList=("cluster01" "cluster02" "cluster03" "cluster04" "cluster05" "cluster06" "cluster07" "cluster08" "cluster09" "cluster10" "cluster11" "cluster12" "cluster13" "cluster14" "cluster15" "cluster16" "cluster17")

start_time=$( date +%s.%N )

#kubectl create -f openmcpdeploy.yaml

cat <<EOF | kubectl create -f -
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: test-deploy2
  namespace: keti
spec:
  replicas: 30
  clusters:
  - cluster01
  - cluster02
  - cluster03
  - cluster04
  - cluster05
  - cluster07
  - cluster09
  - cluster10
  - cluster11
  - cluster12
  - cluster13
  - cluster14
  - cluster15
  - cluster16
  - cluster17
  labels:
      app: openmcp-nginx
  template:
    spec:
      template:
        spec:
          imagePullSecrets:
            - name: regcred
          containers:
            - image: nginx
              name: nginx
              resources:
                requests:
                  memory: "9"
                  cpu: "0.1"
EOF

sleep 1

echo "---"
echo "Wait Until Pod Status is Running ..."


for ((;;))
do
	AllRunningFlag="1"
	for cluster in "${clusterList[@]}"; do
		statuslist=($(kubectl get pod -n keti --context $cluster | grep 'test-deploy' | awk '{print $3}'))

		#echo $statuslist
	        for ((i=0; i<${#statuslist[@]}; i++)); do
	            status=${statuslist[i]}
	            if [ "$status" != "Running" ]; then
	                   AllRunningFlag="0"
	                   break
	            fi
	        done

	        if [ "$AllRunningFlag" == "0" ];
	        then
		   echo "$cluster is Not Ready"
	           break
        	fi
	done

	if [ "$AllRunningFlag" == "1" ];
	then
		break
	fi
done


end_time=$( date +%s.%N --date="$start_time seconds ago" )

for cluster in "${clusterList[@]}"; do
	echo "[$cluster]"
	kubectl get pods -n keti --context $cluster
done

echo "---"
echo "*** 30pods deploy time: ${end_time}s"
