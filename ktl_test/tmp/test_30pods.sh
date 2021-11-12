#!/bin/bash

start_time=$( date +%s.%N )

kubectl create -f test.yaml

echo "---"
echo "Wait Until Pod Status is Running ..."
for ((;;))
do
        tmp=1
        statuslist=($(kubectl get pod -n openmcp | grep 'test-deploy' | awk '{print $3}'))

        for ((i=0; i<${#statuslist[@]}; i++)); do
            status=${statuslist[i]}
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
echo "*** 30pods_deploy_time: ${end_time}s"
echo "---"

kubectl get pods -n openmcp

