#!/bin/bash

start_time=$( date +%s.%N )

cat <<EOF | kubectl create -f -
apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: test-deploy
  namespace: keti
spec:
  replicas: 30
  clusters:
  - cluster01
  - cluster02
  - cluster03
  - cluster04
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
#              resources:
#                requests:
#                  memory: "1"
#                  cpu: "0.1"
EOF

echo "---"
echo "Wait Until Pod Status is Running ..."
for ((;;))
do
        tmp=1

        statuslist0=($(kubectl get pod -n keti --context cluster01 | grep 'test-deploy' | awk '{print $3}'))

        if [ ${#statuslist0[@]} == 0 ]; then
           echo "wait..."
           continue
        fi  

        for ((i=0; i<${#statuslist0[@]}; i++)); do
            status=${statuslist0[i]}
            if [ "$status" != "Running" ]; then
                   tmp=2
                   break
            fi
        done

        if [ "$tmp" == "2" ];
        then
           continue
        fi

        statuslist1=($(kubectl get pod -n keti --context cluster02 | grep 'test-deploy' | awk '{print $3}'))
        
        if [ ${#statuslist1[@]} == 0 ]; then
           echo "wait..."
           continue
        fi

        for ((i=0; i<${#statuslist1[@]}; i++)); do
            status=${statuslist1[i]}
            if [ "$status" != "Running" ]; then
                   tmp=2
                   break
            fi
        done

        if [ "$tmp" == "2" ];
        then
           continue
        fi

        statuslist2=($(kubectl get pod -n keti --context cluster03 | grep 'test-deploy' | awk '{print $3}'))

        if [ ${#statuslist2[@]} == 0 ]; then
           echo "wait..."
           continue
        fi


        for ((i=0; i<${#statuslist2[@]}; i++)); do
            status=${statuslist2[i]}
            if [ "$status" != "Running" ]; then
                   tmp=2
                   break
            fi
        done

        if [ "$tmp" == "2" ];
        then
           continue
        fi

        statuslist3=($(kubectl get pod -n keti --context cluster04 | grep 'test-deploy' | awk '{print $3}'))

        if [ ${#statuslist3[@]} == 0 ]; then
           echo "wait..."
           continue
        fi


        for ((i=0; i<${#statuslist3[@]}; i++)); do
            status=${statuslist3[i]}
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

echo "[cluster1]"
kubectl get pods -n keti --context cluster01
echo "[cluster2]"
kubectl get pods -n keti --context cluster02
echo "[cluster3]"
kubectl get pods -n keti --context cluster03
echo "[cluster4]"
kubectl get pods -n keti --context cluster04

echo "---"
echo "*** 30pods deploy time: ${end_time}s"
