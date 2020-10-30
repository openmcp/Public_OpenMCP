package reference

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

type MyEtcd struct {
	MyCluster *MyCluster
	Endpoints []string
	TLSInfo   transport.TLSInfo
	Client    clientv3.Client
	Cfg       clientv3.Config
}

func NewMyETCD(myCluster *MyCluster) *MyEtcd {
	backup_dirname, _ := filepath.Abs("/etc/pki/etcd-ca")
	openmcp_dirname, _ := filepath.Abs("/home/nfs/openmcp/" + myCluster.OpenMCPMasterIP)

	myetcd := &MyEtcd{
		MyCluster: myCluster,
		Endpoints: []string{myCluster.IP + ":" + myCluster.PORT},
	}

	certFilePath := ""
	keyFilePath := ""
	trustedCAFilePath := ""

	if myCluster.isEtcdBackupServer {
		certFilePath = filepath.Join(backup_dirname, "certs", "etcd-client.crt")
		keyFilePath = filepath.Join(backup_dirname, "private", "etcd-client.key")
		trustedCAFilePath = filepath.Join(backup_dirname, "certs", "ca.crt")
	} else {
		if myCluster.IP == myCluster.OpenMCPMasterIP {
			certFilePath = filepath.Join(openmcp_dirname, "master/pki", "server.crt")
			keyFilePath = filepath.Join(openmcp_dirname, "master/pki", "server.key")
			trustedCAFilePath = filepath.Join(openmcp_dirname, "master/pki", "ca.crt")
		} else {
			certFilePath = filepath.Join(openmcp_dirname, "members/join", myCluster.IP, "pki/server.crt")
			keyFilePath = filepath.Join(openmcp_dirname, "members/join", myCluster.IP, "pki/server.key")
			trustedCAFilePath = filepath.Join(openmcp_dirname, "members/join", myCluster.IP, "pki/ca.crt")
		}
	}
	tlsInfo := transport.TLSInfo{
		CertFile:      certFilePath,
		KeyFile:       keyFilePath,
		TrustedCAFile: trustedCAFilePath,
	}

	myetcd.TLSInfo = tlsInfo

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg := clientv3.Config{
		Endpoints:   myetcd.Endpoints,
		DialTimeout: dialTimeout,
		TLS:         tlsConfig,
	}
	myetcd.Cfg = cfg
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	myetcd.Client = *cli

	return myetcd
}

func (e *MyEtcd) Get(key string) *clientv3.GetResponse {
	resp, _ := e.Client.Get(context.TODO(), key)
	//for _, ev := range resp.Kvs {
	//	key := string(ev.Key)
	//	val := string(ev.Value)
	//	fmt.Println(key, val)
	//}
	return resp
}
func (e *MyEtcd) GetWithPrefix(key string) *clientv3.GetResponse {
	resp, _ := e.Client.Get(context.TODO(), key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	//for _, ev := range resp.Kvs {
	//	key := string(ev.Key)
	//	val := string(ev.Value)
	//	fmt.Println(key, val)
	//}
	return resp
}
func (e *MyEtcd) Put(key, val string) *clientv3.PutResponse {
	resp, err := e.Client.Put(context.TODO(), key, val)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	return resp
}
func (e *MyEtcd) Del(key string) *clientv3.DeleteResponse {
	resp, err := e.Client.Delete(context.TODO(), key)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	return resp

}
func (e *MyEtcd) DelWithPrefix(key string) *clientv3.DeleteResponse {
	resp, err := e.Client.Delete(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	return resp

}
func (e *MyEtcd) WatchPrint(key string) {
	rch := e.Client.Watch(context.TODO(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("[%s] %s %q : %q\n", e.MyCluster.ClusterName, ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}
func (e *MyEtcd) WatchPrintWithPrefix(key string) {
	rch := e.Client.Watch(context.TODO(), key, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("[%s] %s %q : %q\n", e.MyCluster.ClusterName, ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}

func (e *MyEtcd) PutWithMap(etcdMap map[string]string) {
	var wait sync.WaitGroup
	wait.Add(len(etcdMap))
	for k, v := range etcdMap {
		go func(k, v string) {
			defer wait.Done()
			e.Put(k, v)

		}(k, v)
	}

	wait.Wait()
}
