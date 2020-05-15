#!/bin/bash
apt-get update
apt-get install nfs-common nfs-kernel-server portmap

mkdir /home/nfs
chmod 777 /home/nfs

ssh-keygen -f /root/.ssh/id_rsa -t rsa -P ""
mkdir /home/nfs/ssh && cp id_rsa.pub /home/nfs/ssh

mkdir /home/nfs/pv

echo "/home/nfs *(rw,no_root_squash,sync)" >> /etc/exports

systemctl restart nfs-kernel-server
