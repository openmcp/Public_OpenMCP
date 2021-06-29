package etcd

import (
	"context"
	"time"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-snapshot/pkg/util"
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

func InitEtcd() (*Etcd, error) {

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
		omcplog.Error("etcd.go init error : ", err)
		return nil, err
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
	return etcd, nil
}

func (e *Etcd) GetEtcdGroupSnapshot(groupSnapshotKey string) (*clientv3.GetResponse, error) {
	prifix := util.MakePrifix(groupSnapshotKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, getErr := e.Client.Get(ctx, prifix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if getErr != nil {
		omcplog.Error("etcd.go : groupGet Err", getErr)
		return nil, getErr
	}
	return resp, nil
}

func (e *Etcd) Get(key string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, getErr := e.Client.Get(ctx, key)
	cancel()
	if getErr != nil {
		omcplog.Error("etcd.go : get Err", getErr)
		return nil, getErr
	}
	return resp, nil
}

func (e *Etcd) Put(key, val string) (*clientv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, getErr := e.Client.Put(ctx, key, val)
	cancel()
	if getErr != nil {
		omcplog.Error("etcd.go : Put Err", getErr)
		return nil, getErr
	}
	return resp, nil
}
