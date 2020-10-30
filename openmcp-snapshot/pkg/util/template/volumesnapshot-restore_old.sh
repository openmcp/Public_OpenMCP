#!/bin/bash

# 1. Volume 패스의 파일들을 삭제한다. (백업)
cd !FINDPATH
mkdir backup/!DATE
mv !FINDPATH/* backup/!DATE
mv !FINDPATH/.* backup/!DATE


# 2. externalNFS 에서 해당 job 으로 지정된 스냅샷 폴더로 이동한다,
cd /storage    # externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨


# 3. 복구하려는 스냅샷의 압축을 Volume 패스에 푼다. DATE 는 KEY값에서 추출함
export FILES=$(ls -tr)
for FILE in $FILES
do
  tar xvf !DATE -C !FINDPATH
  if [ "$FILE" == "!DATE" ]; then
    break
  fi
done