# 고가용성 클러스터 구축

## 구성 환경

Master1 IP : 10.0.3.100 
Master2 IP : 10.0.3.40  
Master3 IP : 10.0.3.30  
LoadBalancer IP : 10.0.3.99

## 구축 순서
1. hostname 설정
1. docker 설치
1. kubeadm/kubectl/kubelet 설치
1. HAProxy 설치
1. Keepalived 설치
1. 클러스터 생성 (kubeadm init)
1. 마스터 노드 추가 (kubeadm join)
1. 마스터 노드 복제 확인

## 1. hostname 설정
### > Master1
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master1
Please Reconnect shell to Change Hostname!
```
### > Master2
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master2
Please Reconnect shell to Change Hostname!
```
### > Master3
```
$ cd create_cluster
$ ./1.set_host.sh 
Your Hostname? master3
Please Reconnect shell to Change Hostname!
```
## 2. docker 설치
### > Master1 / Master2 / Master3
```
$ apt-get update
$ apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
$ add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
$ apt-cache madison docker-ce
$ apt-get install docker-ce=17.03.2~ce-0~ubuntu-xenial
```
```
$ docker version
Client:
 Version:      17.03.2-ce
 API version:  1.27
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 03:35:14 2017
 OS/Arch:      linux/amd64

Server:
 Version:      17.03.2-ce
 API version:  1.27 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   f5ec1e2
 Built:        Tue Jun 27 03:35:14 2017
 OS/Arch:      linux/amd64
 Experimental: true
```
## 3. kubeadm/kubectl/kubelet 설치
### > Master1 / Master2 / Master3
```
$ apt-get install kubeadm=1.17.3-00
$ apt-get install kubectl=1.17.3-00
$ apt-get install kubelet=1.17.3-00
```
## 4. HAProxy 설치 / config 수정
### > Master1 / Master2 / Master3
```
$ apt-get install haproxy
```
```
$ vi /etc/haproxy/haproxy.cfg
global
        log /dev/log    local0
        log /dev/log    local1 notice
        chroot /var/lib/haproxy
        stats socket /run/haproxy/admin.sock mode 660 level admin
        stats timeout 30s
        user haproxy
        group haproxy
        daemon

        ca-base /etc/ssl/certs
        crt-base /etc/ssl/private
        
        ssl-default-bind-ciphers ECDH+AESGCM:DH+AESGCM:ECDH+AES256:DH+AES256:ECDH+AES128:DH+AES:ECDH+3DES:DH+3DES:RSA+AESGCM:RSA+AES:RSA+3DES:!aNULL:!MD5:!DSS
        ssl-default-bind-options no-sslv3

defaults
        log     global
        mode    http
        option  httplog
        option  dontlognull
        timeout connect 5000
        timeout client  50000
        timeout server  50000
        errorfile 400 /etc/haproxy/errors/400.http
        errorfile 403 /etc/haproxy/errors/403.http
        errorfile 408 /etc/haproxy/errors/408.http
        errorfile 500 /etc/haproxy/errors/500.http
        errorfile 502 /etc/haproxy/errors/502.http
        errorfile 503 /etc/haproxy/errors/503.http
        errorfile 504 /etc/haproxy/errors/504.http

frontend kubernetes-master-lb
bind 10.0.3.99:26443
option tcplog
mode tcp
default_backend kubernetes-master-nodes

backend kubernetes-master-nodes
mode tcp
balance roundrobin
option tcp-check
option tcplog
server node1 10.0.3.100:6443 check
server node2 10.0.3.40:6443 check
server node3 10.0.3.30:6443 check
```
```
$ systemctl enable haproxy
$ systemctl start haproxy
$ systemctl status haproxy
```
## 5. Keepalived 설치 / config 수정
### > Master1 / Master2 / Master3
```
$ apt-get install keepalived
```
### > Master1
```
$ vi /etc/keepalived/keepalived.conf
global_defs {
        smtp_server 127.0.0.1
        smtp_connect_timeout 30
}

vrrp_instance VI_1 {
        state MASTER
        interface eno1
        virtual_router_id 40
        priority 102
        advert_int 1
        authentication {
                auth_type PASS
                auth_pass 1111
        }
        virtual_ipaddress {
                10.0.3.99
        }
}
```
```
$ systemctl enable keepalived
$ systemctl start keepalived
$ systemctl status keepalived
● keepalived.service - Keepalive Daemon (LVS and VRRP)
   Loaded: loaded (/lib/systemd/system/keepalived.service; enabled; vendor preset: enabled)
   Active: active (running) since 목 2020-10-22 14:55:24 KST; 8s ago
  Process: 22752 ExecStart=/usr/sbin/keepalived $DAEMON_ARGS (code=exited, status=0/SUCCESS)
 Main PID: 22758 (keepalived)
    Tasks: 3
   Memory: 2.4M
      CPU: 18ms
   CGroup: /system.slice/keepalived.service
           ├─22758 /usr/sbin/keepalived
           ├─22762 /usr/sbin/keepalived
           └─22764 /usr/sbin/keepalived

10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Registering Kernel netlink command channel
10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Registering gratuitous ARP shared channel
10월 22 14:55:24 master1 Keepalived_healthcheckers[22762]: Opening file '/etc/keepalived/keepalived.conf'.
10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Unable to load ipset library
10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Unable to initialise ipsets
10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Opening file '/etc/keepalived/keepalived.conf'.
10월 22 14:55:24 master1 Keepalived_vrrp[22764]: Using LinkWatch kernel netlink reflector...
10월 22 14:55:24 master1 Keepalived_healthcheckers[22762]: Using LinkWatch kernel netlink reflector...
10월 22 14:55:25 master1 Keepalived_vrrp[22764]: VRRP_Instance(VI_1) Transition to MASTER STATE
10월 22 14:55:26 master1 Keepalived_vrrp[22764]: VRRP_Instance(VI_1) Entering MASTER STATE
```
### > Master2
```
$ vi /etc/keepalived/keepalived.conf
global_defs {
        smtp_server 127.0.0.1
        smtp_connect_timeout 30
}

vrrp_instance VI_2 {
        state BACKUP
        interface enp5s0f0
        virtual_router_id 40
        priority 101
        advert_int 1
        authentication {
                auth_type PASS
                auth_pass 1111
        }
        virtual_ipaddress {
                10.0.3.99
        }
}
```
```
$ systemctl enable keepalived
$ systemctl start keepalived
$ systemctl status keepalived
● keepalived.service - Keepalive Daemon (LVS and VRRP)
   Loaded: loaded (/lib/systemd/system/keepalived.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2020-10-21 22:55:25 PDT; 8s ago
  Process: 17302 ExecStart=/usr/sbin/keepalived $DAEMON_ARGS (code=exited, status=0/SUCCESS)
 Main PID: 17305 (keepalived)
    Tasks: 3
   Memory: 1.0M
      CPU: 16ms
   CGroup: /system.slice/keepalived.service
           ├─17305 /usr/sbin/keepalived
           ├─17306 /usr/sbin/keepalived
           └─17307 /usr/sbin/keepalived

Oct 21 22:55:25 master2 Keepalived_healthcheckers[17306]: Registering Kernel netlink command channel
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Registering Kernel netlink command channel
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Registering gratuitous ARP shared channel
Oct 21 22:55:25 master2 Keepalived_healthcheckers[17306]: Opening file '/etc/keepalived/keepalived.conf'.
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Unable to load ipset library
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Unable to initialise ipsets
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Opening file '/etc/keepalived/keepalived.conf'.
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: Using LinkWatch kernel netlink reflector...
Oct 21 22:55:25 master2 Keepalived_vrrp[17307]: VRRP_Instance(VI_2) Entering BACKUP STATE
```
### > Master3
```
$ vi /etc/keepalived/keepalived.conf
global_defs {
        smtp_server 127.0.0.1
        smtp_connect_timeout 30
}

vrrp_instance VI_3 {
        state BACKUP
        interface enp96s0f0
        virtual_router_id 40
        priority 100
        advert_int 1
        authentication {
                auth_type PASS
                auth_pass 1111
        }
        virtual_ipaddress {
                10.0.3.99
        }
}
```
```
$ systemctl enable keepalived
$ systemctl start keepalived
$ systemctl status keepalived
● keepalived.service - Keepalive Daemon (LVS and VRRP)
   Loaded: loaded (/lib/systemd/system/keepalived.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2020-10-21 22:55:51 PDT; 1s ago
  Process: 73064 ExecStart=/usr/sbin/keepalived $DAEMON_ARGS (code=exited, status=0/SUCCESS)
 Main PID: 73066 (keepalived)
    Tasks: 3
   Memory: 1.4M
      CPU: 22ms
   CGroup: /system.slice/keepalived.service
           ├─73066 /usr/sbin/keepalived
           ├─73067 /usr/sbin/keepalived
           └─73068 /usr/sbin/keepalived

Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Registering Kernel netlink reflector
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Registering Kernel netlink command channel
Oct 21 22:55:51 master3 Keepalived_healthcheckers[73067]: Opening file '/etc/keepalived/keepalived.conf'.
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Registering gratuitous ARP shared channel
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Unable to load ipset library
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Unable to initialise ipsets
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Opening file '/etc/keepalived/keepalived.conf'.
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: Using LinkWatch kernel netlink reflector...
Oct 21 22:55:51 master3 Keepalived_vrrp[73068]: VRRP_Instance(VI_3) Entering BACKUP STATE
```

## 6. 클러스터 생성 (kubeadm init)
### > Master1
```
$ kubeadm init --pod-network-cidr=10.244.0.0/16 --ignore-preflight-errors=NumCPU --v=5 --control-plane-endpoint "10.0.3.99:26443" --upload-certs
...
...
...

Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of the control-plane node running the following command on each as root:

  kubeadm join 10.0.3.99:26443 --token ulr0v8.vsz2jwjk5uvvl7h1 \
    --discovery-token-ca-cert-hash sha256:1c15018f135af60b8b44915cf3e5030dbdb14e5f56a8365c9cd3a97c5a776e31 \
    --control-plane --certificate-key 07a8671ee05aae5961d1b220ff25059c2b6fd4896de8c940099921b0b680c770

Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use
"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 10.0.3.99:26443 --token ulr0v8.vsz2jwjk5uvvl7h1 \
    --discovery-token-ca-cert-hash sha256:1c15018f135af60b8b44915cf3e5030dbdb14e5f56a8365c9cd3a97c5a776e31 
```
```
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
$ sysctl net.bridge.bridge-nf-call-iptables=1
$ kubectl apply -f calico.yaml
$ kubectl get pods -A
NAMESPACE     NAME                                       READY   STATUS             RESTARTS   AGE
kube-system   calico-kube-controllers-8464785d6b-8f2gm   1/1     Running            0          103s
kube-system   calico-node-287sb                          1/1     Running            0          103s
kube-system   coredns-6955765f44-b5mx4                   1/1     Running            0          2m8s
kube-system   coredns-6955765f44-hzncd                   1/1     Running            0          2m8s
kube-system   etcd-master1                               1/1     Running            0          2m17s
kube-system   kube-apiserver-master1                     1/1     Running            0          2m17s
kube-system   kube-controller-manager-master1            1/1     Running            0          2m17s
kube-system   kube-proxy-6rc5j                           1/1     Running            0          2m8s
kube-system   kube-scheduler-master1                     1/1     Running            0          2m17s
```
## 7. 마스터 노드 추가 (kubeadm join)
### > Master2
```
$ kubeadm join 10.0.3.99:26443 --token ulr0v8.vsz2jwjk5uvvl7h1 \
    --discovery-token-ca-cert-hash sha256:1c15018f135af60b8b44915cf3e5030dbdb14e5f56a8365c9cd3a97c5a776e31 \
    --control-plane --certificate-key 07a8671ee05aae5961d1b220ff25059c2b6fd4896de8c940099921b0b680c770
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```
### > Master3
```
$ kubeadm join 10.0.3.99:26443 --token ulr0v8.vsz2jwjk5uvvl7h1 \
    --discovery-token-ca-cert-hash sha256:1c15018f135af60b8b44915cf3e5030dbdb14e5f56a8365c9cd3a97c5a776e31 \
    --control-plane --certificate-key 07a8671ee05aae5961d1b220ff25059c2b6fd4896de8c940099921b0b680c770
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

## 8. 마스터 노드 복제 확인
### > Master1 / Master2 / Master3
```
$ kubectl get nodes
NAME      STATUS   ROLES    AGE     VERSION
master1   Ready    master   3m15s   v1.17.3
master2   Ready    master   4m6s    v1.17.3
master3   Ready    master   2m13s   v1.17.3

$ kubectl get pods -A
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-8464785d6b-8f2gm   1/1     Running   0          21m
kube-system   calico-node-287sb                          1/1     Running   0          21m
kube-system   calico-node-tjp8n                          1/1     Running   0          18m
kube-system   calico-node-zdw2z                          1/1     Running   0          18m
kube-system   coredns-6955765f44-dlzl7                   1/1     Running   0          19m
kube-system   coredns-6955765f44-dxq4n                   1/1     Running   0          22m
kube-system   etcd-master1                               1/1     Running   0          22m
kube-system   etcd-master2                               1/1     Running   0          18m
kube-system   etcd-master3                               1/1     Running   0          18m
kube-system   kube-apiserver-master1                     1/1     Running   0          22m
kube-system   kube-apiserver-master2                     1/1     Running   0          18m
kube-system   kube-apiserver-master3                     1/1     Running   0          18m
kube-system   kube-controller-manager-master1            1/1     Running   0          22m
kube-system   kube-controller-manager-master2            1/1     Running   0          18m
kube-system   kube-controller-manager-master3            1/1     Running   0          18m
kube-system   kube-proxy-5frhq                           1/1     Running   0          18m
kube-system   kube-proxy-6rc5j                           1/1     Running   0          21m
kube-system   kube-proxy-cchrl                           1/1     Running   0          18m
kube-system   kube-scheduler-master1                     1/1     Running   0          22m
kube-system   kube-scheduler-master2                     1/1     Running   0          18m
kube-system   kube-scheduler-master3                     1/1     Running   0          18m
```
