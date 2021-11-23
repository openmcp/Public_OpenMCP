#!/bin/bash
set -e
# 1. externalNFS 에서 해당 deploy 로 지정된 스냅샷 폴더로 이동한다,
#sleep 1000000000000000000
# storage - externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨
cd /storage/!PATH

ls -al --time-style="+%Y%m%d-%H%M%S" | grep -v ^d | grep -v total |grep -v backup  | awk '
    BEGIN { ORS = ""; print " [ "}
    { printf "%s{\"size\": \"%s\", \"snapshotKey\": \"%s\", \"date\": \"%s\", \"pvName\": \"!PVNAME\"}",
          separator, $5, $7, $6
      separator = ", "
    }
    END { print " ] " }
'

#  [ {"size": "10240", "snapshotKey": "1636449053"}, {"size": "10240", "snapshotKey": "1637027140"} ] 

#아래 명령어가 없으면 job 이 완료가 아닌 running 상태가 됨.
touch /success
