package manager

import (
	"context"
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/snapshot"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/pkg/transport"
	"go.uber.org/zap"
	"log"
	"path/filepath"
	"time"
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
		Endpoints: []string{myCluster.IP + ":2379"},
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
func (e *MyEtcd) SnapShot(filename string) (string, error) {

	if filename == "" {
		filename = fmt.Sprintf("snapshot%d.db", time.Now().Unix())
	}
	ccfg := e.Cfg

	sp := snapshot.NewV3(zap.NewExample())
	path, err := filepath.Abs(filepath.Join(SNAPSHOT_DIR, e.MyCluster.ClusterName))

	err = ensureDir(path)

	dpPath := filepath.Join(path, filename)
	if err = sp.Save(context.Background(), etcd.Config(ccfg), dpPath); err != nil {
		fmt.Println(err)
	}
	return dpPath, err
}

func (e *MyEtcd) Restore(db_filepath string) string {
	const testClusterTkn = "tkn"

	snapshotDir := SNAPSHOT_DIR

	//cURLs, pURLs := urls[:clusterN], urls[clusterN:]

	cfg := embed.NewConfig()
	cfg.Name = e.MyCluster.ClusterName
	cfg.InitialClusterToken = testClusterTkn
	cfg.InitialCluster = fmt.Sprintf("%s=%s", cfg.Name, "http://"+e.MyCluster.IP+":2379")
	cfg.Dir = filepath.Join(snapshotDir, e.MyCluster.ClusterName, fmt.Sprintf("%d", time.Now().Unix()))

	sp := snapshot.NewV3(zap.NewExample())

	pss := make([]string, 0, 1)
	pss = append(pss, "http://"+e.MyCluster.IP+":2379")

	if err := sp.Restore(snapshot.RestoreConfig{
		SnapshotPath:        db_filepath,
		Name:                cfg.Name,
		OutputDataDir:       cfg.Dir,
		InitialCluster:      cfg.InitialCluster,
		InitialClusterToken: cfg.InitialClusterToken,
		PeerURLs:            pss,
	}); err != nil {
		fmt.Println(err)
	}
	return cfg.Dir

	//cmdStr := "ssh "+e.MyCluster.IP+" rm -rf /var/lib/etcd"
	//_, err = CmdExec(cmdStr)
	//if err != nil {
	//	fmt.Println("Err !", err)
	//}
	//cmdStr = "scp -r "+ cfg.Dir + " root@"+e.MyCluster.IP+":/var/lib/etcd"
	//_, err = CmdExec(cmdStr)
	//if err != nil {
	//	fmt.Println("Err !", err)
	//}
	//cmdStr = "systemctl stop etcd.service"
	//_, err = CmdExec(cmdStr)
	//if err != nil {
	//	fmt.Println("Err !", err)
	//}
	//cmdStr = "systemctl start etcd.service"
	//_, err = CmdExec(cmdStr)
	//if err != nil {
	//	fmt.Println("Err !", err)
	//}
	//os.RemoveAll(cfg.Dir)

}
