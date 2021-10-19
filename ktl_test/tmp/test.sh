#!/bin/bash

start_time=$( date +%s.%N )


echo "---"
echo "Wait Until Pod Status is Running ..."
for ((;;))
do
        tmp=1
        statuslist=($(kubectl get pod -n keti --context cluster1 | grep 'test-deploy' | awk '{print $3}'))

        for ((i=0; i<${#statuslist[@]}; i++)); do
            status=${statuslist[i]}
            echo $status
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
            echo $status
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

echo "---"
echo "*** 30pods deploy time: ${end_time}s"
