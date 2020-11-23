#!/bin/bash
# 1. externalNFS 에서 해당 deploy 로 지정된 스냅샷 폴더로 이동한다,

cd /storage    # externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨
mkdir -p !PATH
cd !PATH

export lastDir=`ls -tr | tail -1`  #가장 최근 스냅샷 폴더

# 2. newerthan, 을 구한다. 폴더가 비어있을 경우 newerthan 는 1970년1월1일이다.
echo "newerthan init"
date +"%F%Z%T" --date=@0
export newerthan=`date +"%F%Z%T" --date=@0`  #초기화
if [ -n "$lastDir" ]; then
  echo "newerthan mod"
  date +"%F%Z%T" --date=@$lastDir
  export newerthan=`date +"%F%Z%T" --date=@$lastDir`   #가장 최근에 스냅샷한 시간
fi
echo "olderthan mod"
date +"%F%Z%T" --date=@!DATE
# 3. olderthan 을 구한다. olderthan 은 현재 시간(리눅스시간) 이다. -> 이것은 코드상에서 계산에서 넣도록한다.
export olderthan=`date +"%F%Z%T" --date=@!DATE`      # 스냅샷 시작 시간

echo "list"
find /data -type f -newermt "$newerthan" ! -newermt "$olderthan"

# 4. newerthan, olderthan 을 이용하여 파일 찾아서 압축   #/data 인 이유는 PV 에 연결된 /data가 여깃음
find /data -type f -newermt "$newerthan" ! -newermt "$olderthan" | xargs tar cvf !DATE

echo "Snapshot end"