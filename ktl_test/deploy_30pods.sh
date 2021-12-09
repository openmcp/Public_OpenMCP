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
              imagePullPolicy: IfNotPresent
              name: nginx
EOF

echo "---"
echo "Wait Until Pod Status is Running ..."
for ((;;))
do
        tmp=1

<<<<<<< HEAD
        statuslist0=($(kubectl get pod -n keti --context cluster01 | grep 'test-deploy' | awk '{print $3}'))

        if [ ${#statuslist0[@]} == 0 ]; then
           echo "wait..."
=======
        statuslist0=$(kubectl get pod -n keti --context cluster01 2>&1 | grep 'test-deploy' | awk '{print $3}')

        if [ "${statuslist0}" == "" ]; then
>>>>>>> ff3bfe7f8885714229f45760f6f759adf6980acf
           continue
        else
           statuslist0=($(kubectl get pod -n keti --context cluster01 | grep 'test-deploy' | awk '{print $3}'))
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

<<<<<<< HEAD
        statuslist1=($(kubectl get pod -n keti --context cluster02 | grep 'test-deploy' | awk '{print $3}'))
        
        if [ ${#statuslist1[@]} == 0 ]; then
           echo "wait..."
=======
        statuslist1=$(kubectl get pod -n keti --context cluster02 2>&1 | grep 'test-deploy' | awk '{print $3}')

        if [ "${statuslist1}" == "" ]; then
>>>>>>> ff3bfe7f8885714229f45760f6f759adf6980acf
           continue
        else
           statuslist1=($(kubectl get pod -n keti --context cluster02 | grep 'test-deploy' | awk '{print $3}'))
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


        statuslist2=$(kubectl get pod -n keti --context cluster03 2>&1 | grep 'test-deploy' | awk '{print $3}')


        if [ "${statuslist2}" == "" ]; then
           continue
        else
           statuslist2=($(kubectl get pod -n keti --context cluster03 | grep 'test-deploy' | awk '{print $3}'))

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

        statuslist3=$(kubectl get pod -n keti --context cluster04 2>&1 | grep 'test-deploy' | awk '{print $3}')

        if [ "${statuslist3}" == "" ]; then
           continue
        else
           statuslist3=($(kubectl get pod -n keti --context cluster04 | grep 'test-deploy' | awk '{print $3}'))
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
