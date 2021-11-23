#!/bin/bash
set -e
echo "1. Info, Backup"
echo "PATH(ExternalNFS) : !PATH"
echo "DATE : !DATE"

# 1. Volume 패스의 파일들을 삭제한다. (백업)
echo "@@ backup"
#sleep 1000000
cd /data
mkdir -p /storage/!PATH/backup/!DATE
echo "add Data : /storage/!PATH/backup/!DATE"
cp -r /data/ /storage/!PATH/backup/!DATE
#cp -r /data/.* /storage/!PATH/backup/!DATE

echo "@@ list Data : /storage/!PATH/backup/!DATE"
ls /storage/!PATH/backup/!DATE

# 2. externalNFS 에서 해당 job 으로 지정된 스냅샷 폴더로 이동한다,
# /storage : externalNFS 의 /home/nfs/storage/CLUSTERNAME/volume/PVNAME/ 와 마운트됨
#mkdir -p .!PATH
echo "2. move ExternalNFS folder"
echo "cd /storage/!PATH"
cd /storage/!PATH



echo "3. Unzips..."
tar xfP !DATE --listed-incremental backuplist 



## 3. 복구하려는 스냅샷의 압축을 Volume 패스에 푼다. DATE 는 KEY값에서 추출함
#echo "3. Unzips..."
#export FILES="$(ls -tr)"
#echo "Files : "
#echo "$FILES"
#
#if [ -z "$FILES" ]; then
#  echo "target File empty!"
#  touch /success
#  exit 100
#fi
#for FILE in $FILES
#do
#  if [ "$FILE" == "backup" ]; then
#    echo "End backup folder"
#    break
#  fi
#  
#  total=`expr "${FILE}" "-" "!DATE"`
#  if [ ${total} -gt 0 ]; then
#    echo "End !DATE tar"
#    break
#  fi
#
#  #echo "unzip... $FILE -C /"
#  tar xfP ${FILE} -C /
#  
#  if [ "$FILE" == "!DATE" ]; then
#    echo "End !DATE tar"
#    break
#  fi
#done
echo "4. Snapshot restore end"
#아래 명령어가 없으면 job 이 완료가 아닌 running 상태가 됨.
touch /success