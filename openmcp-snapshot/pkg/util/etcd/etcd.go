package etcd

import (
	"context"
	"log"
	"time"

	config "openmcp/openmcp/openmcp-snapshot/pkg/util"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

const openmcpDirnam = "/home/nfs/openmcp/" + config.MASTER_IP + "/"

type etcdInfoParamsMap struct {
	Endpoints      []string
	TLSInfo        transport.TLSInfo
	RequestTimeout time.Duration
	DialTimeout    time.Duration
}

// EtcdInfo 는 ETCD 접속정보 입니다.
var EtcdInfo = etcdInfoParamsMap{
	DialTimeout: 2 * time.Second,
	Endpoints:   []string{config.EXTERNAL_ETCD},
	TLSInfo: transport.TLSInfo{
		CertFile:      openmcpDirnam + "certs/etcd-client.crt",
		KeyFile:       openmcpDirnam + "privateetcd-client.key",
		TrustedCAFile: openmcpDirnam + "certs/ca.crt",
	},
	RequestTimeout: 10 * time.Second,
}

type Etcd struct {
	Endpoints      []string
	TLSInfo        transport.TLSInfo
	Client         clientv3.Client
	Cfg            clientv3.Config
	RequestTimeout time.Duration
	DialTimeout    time.Duration
}

func InitEtcd() *Etcd {

	etcd := &Etcd{
		Endpoints:      EtcdInfo.Endpoints,
		TLSInfo:        EtcdInfo.TLSInfo,
		RequestTimeout: EtcdInfo.RequestTimeout,
		DialTimeout:    EtcdInfo.DialTimeout,
	}

	//Init 에 따로 빼둘 필요성 있음.
	//tlsConfig, err := etcd.TLSInfo.ClientConfig()
	//if err != nil {
	//	log.Fatal(err)
	//}
	cfg := clientv3.Config{
		Endpoints:   etcd.Endpoints,
		DialTimeout: etcd.DialTimeout,
		//TLS:         tlsConfig,
	}
	etcd.Cfg = cfg
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	etcd.Client = *cli

	//	ctx, _ := context.WithTimeout(context.Background(), util.RequestTimeout)
	//	cli, err := clientv3.New(clientv3.Config{
	//		DialTimeout: EtcdInfo.DialTimeout,
	//		Endpoints:   EtcdInfo.Endpoints,
	//		TLS:         tlsConfig,
	//	})
	//	if err != nil {
	//		// handle error!
	//		fmt.Println(err)
	//		return nil, err
	//	}
	return etcd
}

func (e *Etcd) Get(key string) *clientv3.GetResponse {
	resp, _ := e.Client.Get(context.TODO(), key)
	//for _, ev := range resp.Kvs {
	//	key := string(ev.Key)
	//	val := string(ev.Value)
	//	fmt.Println(key, val)
	//}
	return resp
}

func (e *Etcd) Put(key, val string) *clientv3.PutResponse {
	resp, err := e.Client.Put(context.TODO(), key, val)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp)
	return resp
}
