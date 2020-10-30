package etcd

import (
	"context"
	"fmt"
	"strconv"

	"openmcp/openmcp/openmcp-snapshot/pkg/util"

	"go.etcd.io/etcd/clientv3"
)

// InsertEtcd : 키 추가
func InsertEtcd(key string, value string) (bool, error) {

	ctx, _ := context.WithTimeout(context.Background(), util.RequestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: util.EtcdInfo.DialTimeout,
		Endpoints:   util.EtcdInfo.Endpoints,
		//TLS:         tlsConfig,
	})
	if err != nil {
		// handle error!
		fmt.Println(err)
		return false, err
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	//==================================================

	fmt.Println("*** Insert()")
	// Delete all keys ("key prefix")
	//kv.Delete(ctx, "key", clientv3.WithPrefix())
	fmt.Println("key: " + key)
	if value == "" {
		fmt.Println("Value: is nil")
		return false, fmt.Errorf("Value is empty")
	}

	// Insert a key value
	pr, err := kv.Put(ctx, key, value)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	rev := pr.Header.Revision
	fmt.Println("Insert Revision:", rev)

	return true, nil
}

// GetEtcd : 키 검색
func GetEtcd(key string) (string, error) {

	ctx, _ := context.WithTimeout(context.Background(), util.RequestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: util.EtcdInfo.DialTimeout,
		Endpoints:   util.EtcdInfo.Endpoints,
		//TLS:         tlsConfig,
	})
	if err != nil {
		// handle error!
		fmt.Println(err)
		return "", err
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	//==================================================

	fmt.Println("*** GetEtcd()")
	// Delete all keys ("key prefix")
	//kv.Delete(ctx, "key", clientv3.WithPrefix())
	fmt.Println("key: " + key)

	gr, err := kv.Get(ctx, key)
	if err != nil {
		// handle error!
		fmt.Println(err)
		return "", err
	}
	//fmt.Println(gr)

	if gr.Kvs == nil {
		fmt.Println("Value: is nil")
		return "", fmt.Errorf("Value is empty")
	}

	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	return string(gr.Kvs[0].Value), nil
}

// InsertMultiEtcd : 키 추가 멀티
func InsertMultiEtcd(ctx context.Context, kv clientv3.KV) {
	fmt.Println("*** InsertMulti()")
	// Delete all keys ("key prefix")
	//kv.Delete(ctx, "key", clientv3.WithPrefix())

	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)
		kv.Put(ctx, k, strconv.Itoa(i))
	}

	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(3),
	}

	gr, _ := kv.Get(ctx, "key", opts...)

	fmt.Println("--- First page ---")
	for _, item := range gr.Kvs {
		fmt.Println(string(item.Key), string(item.Value))
	}

	lastKey := string(gr.Kvs[len(gr.Kvs)-1].Key)

	fmt.Println("--- Second page ---")
	opts = append(opts, clientv3.WithFromKey())
	gr, _ = kv.Get(ctx, lastKey, opts...)

	// Skipping the first item, which the last item from from the previous Get
	for _, item := range gr.Kvs[1:] {
		fmt.Println(string(item.Key), string(item.Value))
	}
}
