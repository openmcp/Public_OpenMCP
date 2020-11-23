#!/bin/bash

# 1. Volume 패스의 파일들을 삭제한다. (백업)
#sleep 1000000
cd /data
mkdir -p /storage/!PATH/backup/!DATE
echo "add Data : /storage/!PATH/backup/!DATE"
cp -r /data/* /storage/!PATH/backup/!DATE
#cp -r /data/.* /storage/!PATH/backup/!DATE


echo "list Data : /storage/!PATH/backup/!DATE list"
ls /storage/!PATH/backup/!DATE

# 2. externalNFS 에서 해당 job 으로 지정된 스냅샷 폴더로 이동한다,
# /storage : externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨
#mkdir -p .!PATH
echo "cd /storage/!PATH"
cd /storage/!PATH

# 3. 복구하려는 스냅샷의 압축을 Volume 패스에 푼다. DATE 는 KEY값에서 추출함
export FILES="$(ls -tr)"
echo "Files : "
echo "$FILES"
for FILE in $FILES
do
  echo "tar xvf !DATE -C /"
  tar xvf !DATE -C /
  if [ "$FILE" == "!DATE" ]; then
    echo "End !DATE tar"
    break
  fi
done
echo "Snapshot restore end"