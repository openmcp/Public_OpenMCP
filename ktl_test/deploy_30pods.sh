#!/bin/bash

start_time=$( date +%s.%N )

#kubectl create -f openmcpdeploy.yaml

cat <<EOF | kubectl apply -f -
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: test-deploy
  namespace: keti
spec:
  replicas: 3
  clusters:
  - cluster14
  - cluster13
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
        tmp=1
        statuslist=($(kubectl get pod -n keti --context cluster1 | grep 'test-deploy' | awk '{print $3}'))

        for ((i=0; i<${#statuslist[@]}; i++)); do
            status=${statuslist[i]}
            if [ "$status" != "Running" ]; then
                   tmp=2
                   break
            fi
        done

        if [ "$tmp" == "2" ];
        then
           continue
        fi

        statuslist2=($(kubectl get pod -n keti --context cluster2 | grep 'test-deploy' | awk '{print $3}'))

        for ((i=0; i<${#statuslist2[@]}; i++)); do
            status=${statuslist2[i]}
            if [ "$status" != "Running" ]; then
                   tmp=2
                   break
            fi
        done

        if [ "$tmp" == "1" ];
        then
           break
        fi
done

end_time=$( date +%s.%N --date="$start_time seconds ago" )

kubectl get pods -n keti --context cluster1
kubectl get pods -n keti --context cluster2
echo "---"
echo "*** 30pods deploy time: ${end_time}s"
