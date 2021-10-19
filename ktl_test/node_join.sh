#!/bin/bash

start_time=$( date +%s.%N )

#ssh root@$1 "kubectl get pods -A"
ssh root@$2 "kubectl node-size $3"

sleep 1
#kubectl node-size $1

echo "---"
echo "Wait Until Node Status is Ready ..."
for ((;;))
do
        tmp=1
        statuslist=($(kubectl get node --context $1 | awk '{print $2}'))
        for ((i=1; i<${#statuslist[@]}; i++)); do
            status=${statuslist[i]}
            if [ "$status" != "Ready" ]; then
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

kubectl get node --context cluster2

echo "---"
echo "*** node join time: ${end_time}s"
