package manager

import (
	"context"
	"fmt"
	//
	//"github.com/coreos/etcd/pkg/wait"
	"sync"

	//"fmt"

	//"fmt"

	//"fmt"
	//etcd "github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3"
	//"go.etcd.io/etcd/clientv3/snapshot"
	//"go.etcd.io/etcd/embed"
	//"go.uber.org/zap"
	"io/ioutil"
	"log"
	//"os"
	"path/filepath"
	//"runtime"
	//"sync"
	//"time"
)
type EtcdManager struct {
	MyEtcdBackup *MyEtcd
	MyEtcdClone *MyEtcd
	MyEtcdList []*MyEtcd
	MyEtcdMap map[string]*MyEtcd
}
func NewEtcdManager(OpenMCP_IP, Etcd_IP string) *EtcdManager{

	myClusterEtcd := &MyCluster{"backup_server", Etcd_IP,"2379","", true}
	myEtcdBackup := NewMyETCD(myClusterEtcd)

	myClusterEtcd_Clone := &MyCluster{"clone_server", Etcd_IP, "2369", "", true}
	myEtcdClone := NewMyETCD(myClusterEtcd_Clone)

	myEtcdList := []*MyEtcd{}
	myEtcdMap := make(map[string]*MyEtcd)

	// OpenMCP Master
	filename, _ := filepath.Abs("/home/nfs/openmcp/"+OpenMCP_IP+"/master/config/config")
	yaml, _ := LoadYAMLFile(filepath.Join(filename))
	cluster_name := yaml.Clusters[0].Name
	IP := SplitAny(yaml.Clusters[0].Cluster.Server, "/:")[1]

	myCluster :=  &MyCluster{cluster_name, IP,"2379",OpenMCP_IP, false}

	myEtcd := NewMyETCD(myCluster)

	myEtcdList = append(myEtcdList, myEtcd)
	myEtcdMap[myCluster.ClusterName] = myEtcd

	// Member Clusters
	dirname, _ := filepath.Abs("/home/nfs/openmcp/"+OpenMCP_IP+"/members/join")
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for _, member := range dirs {
		membername := member.Name()
		yaml, _ := LoadYAMLFile(filepath.Join(dirname,membername,"config/config"))

		cluster_name := yaml.Clusters[0].Name
		IP := SplitAny(yaml.Clusters[0].Cluster.Server, "/:")[1]

		myCluster :=  &MyCluster{cluster_name, IP,"2379",OpenMCP_IP, false}

		myEtcd := NewMyETCD(myCluster)

		myEtcdList = append(myEtcdList, myEtcd)
		myEtcdMap[myCluster.ClusterName] = myEtcd
	}

	etcdManager := &EtcdManager{
		MyEtcdBackup: myEtcdBackup,
		MyEtcdClone: myEtcdClone,
		MyEtcdList: myEtcdList,
		MyEtcdMap: myEtcdMap,
	}
	return etcdManager
}

func (e *EtcdManager) MyFuncSyncAll(){
	//runtime.GOMAXPROCS(16)

	key_num := 0
	cluster_db_map := make(map[string]string)
	for cluster_name, myEtcd := range e.MyEtcdMap {
		resp, _ := myEtcd.Client.Get(context.TODO(), "/", clientv3.WithPrefix())
		//key_num = key_num + len(resp.Kvs)
		for _, ev := range resp.Kvs {
			key := "/"+cluster_name+string(ev.Key)
			val := string(ev.Value)
			cluster_db_map[key] = val
		}

	}
	backup_etcd_db_map := make(map[string]string)
	resp, _ := e.MyEtcdBackup.Client.Get(context.TODO(), "/", clientv3.WithPrefix())
	//key_num = key_num + len(resp.Kvs)
	for _, ev := range resp.Kvs {
		key := string(ev.Key)
		val := string(ev.Value)
		backup_etcd_db_map[key] = val
	}

	both_keys, onlyA_keys, onlyB_keys := MapsKeyCompare(backup_etcd_db_map, cluster_db_map)

	key_num = len(both_keys) + len(onlyA_keys) + len(onlyB_keys)

	var wait sync.WaitGroup
	wait.Add(key_num)

	for _, key := range both_keys {
		val := cluster_db_map[key]
		go func(key, val string){
			defer wait.Done()
			_, error := e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
			if error != nil{
				for error != nil {
					fmt.Println(error)
					_, error = e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
				}
			}
			fmt.Println("PUT "+ " : ", key)
		}(key, val)
	}

	for _, key := range onlyA_keys {
		go func(key string){
			defer wait.Done()
			_, error := e.MyEtcdBackup.Client.Delete(context.TODO(), key)
			if error != nil{
				for error != nil {
					fmt.Println(error)
					_, error = e.MyEtcdBackup.Client.Delete(context.TODO(), key)
				}
			}
			fmt.Println("Del "+ " : ", key)
		}(key)
	}
	for _, key := range onlyB_keys {
		val := cluster_db_map[key]
		go func(key, val string){
			defer wait.Done()
			_, error := e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
			if error != nil{
				for error != nil {
					fmt.Println(error)
					_, error = e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
				}
			}
			fmt.Println("PUT "+ " : ", key)
		}(key, val)
	}

	wait.Wait()

}

//func (e *EtcdManager) MyFuncSyncAll(){
//	runtime.GOMAXPROCS(16)
//
//	key_num := 0
//	cluster_db_map := make(map[string]string)
//	for cluster_name, myEtcd := range e.MyEtcdMap {
//		resp, _ := myEtcd.Client.Get(context.TODO(), "/", clientv3.WithPrefix())
//		key_num = key_num + len(resp.Kvs)
//		for _, ev := range resp.Kvs {
//			key := "/"+cluster_name+string(ev.Key)
//			val := string(ev.Value)
//			cluster_db_map[key] = val
//		}
//
//	}
//	var wait sync.WaitGroup
//	wait.Add(key_num)
//	for key, val := range cluster_db_map {
//		go func(key, val string){
//			defer wait.Done()
//			_, error := e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
//			if error != nil{
//				for error != nil {
//					fmt.Println(error)
//					_, error = e.MyEtcdBackup.Client.Put(context.TODO(), key, val)
//				}
//			}
//			//fmt.Println("PUT "+ " : ", key)
//		}(key, val)
//	}
//	wait.Wait()
//
//}
//func (e *EtcdManager) MyFuncRestore(key string, cluster_name string){
//	resp, _ := e.MyEtcdBackup.Client.Get(context.Background(), key)
//	for _, ev := range resp.Kvs {
//		k := string(ev.Key)
//		v := string(ev.Value)
//		e.MyEtcdMap[cluster_name].Client.Put(context.Background(), k, v)
//		fmt.Println(k, v)
//	}
//}
//
//func (e *EtcdManager) MyFuncSnapShot(){
//
//	ccfg := e.MyEtcdBackup.Cfg
//
//	sp := snapshot.NewV3(zap.NewExample())
//	path, err := filepath.Abs("./snapshot")
//
//	dpPath := filepath.Join(path, fmt.Sprintf("snapshot%d.db", time.Now().Unix()))
//	fmt.Println(dpPath)
//	if err = sp.Save(context.Background(), etcd.Config(ccfg), dpPath); err != nil {
//		fmt.Println(err)
//	}
//
//}
//const testClusterTkn = "tkn"
//func (e *EtcdManager) MyFuncRestoreWithSnapshot(filename string) string{
//
//	etcd_dir_path,_ := filepath.Abs("./snapshot")
//	dbPath, _ := filepath.Abs(filepath.Join(etcd_dir_path, filename))
//
//	clusterN := 1
//	urls := newEmbedURLs(clusterN * 2)
//	cURLs, pURLs := urls[:clusterN], urls[clusterN:]
//
//	cfg := embed.NewConfig()
//	cfg.Name = "s1"
//	cfg.InitialClusterToken = testClusterTkn
//	cfg.ClusterState = "existing"
//	cfg.LCUrls, cfg.ACUrls = cURLs, cURLs
//	cfg.LPUrls, cfg.APUrls = pURLs, pURLs
//	cfg.InitialCluster = fmt.Sprintf("%s=%s", cfg.Name, pURLs[0].String())
//	cfg.Dir = filepath.Join(etcd_dir_path, fmt.Sprintf("%d", time.Now().Unix()))
//
//	//cfg.Dir = "/var/lib/etcd"
//
//	sp := snapshot.NewV3(zap.NewExample())
//
//	pss := make([]string, 0, len(pURLs))
//	for _, p := range pURLs {
//		pss = append(pss, p.String())
//	}
//	if err := sp.Restore(snapshot.RestoreConfig{
//		SnapshotPath:        dbPath,
//		Name:                cfg.Name,
//		OutputDataDir:       cfg.Dir,
//		InitialCluster:      cfg.InitialCluster,
//		InitialClusterToken: cfg.InitialClusterToken,
//		PeerURLs:            pss,
//	}); err != nil {
//		fmt.Println(err)
//	}
//	return cfg.Dir

	//cmdStr := "rm -rf /var/lib/etcd"
	//_, err := CmdExec(cmdStr)
	//if err != nil {
	//	fmt.Println("Err !", err)
	//}
	//cmdStr = "cp -r "+ cfg.Dir+" /var/lib/etcd"
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

	//srv, err := embed.StartEtcd(cfg)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer func() {
	//	//os.RemoveAll(cfg.Dir)
	//	//srv.Close()
	//}()
	//select {
	//case <-srv.Server.ReadyNotify():
	//case <-time.After(3 * time.Second):
	//	fmt.Println("failed to start restored etcd member")
	//}

	//resp ,_ := e.myEtcdBackup.Client.Get(context.TODO(), "/", clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	//for _, ev := range resp.Kvs {
	//	key := string(ev.Key)
	//	val := string(ev.Value)
	//	fmt.Println(key, val)
	//}

//}

//func (e *EtcdManager) MyFuncRestoreCluster(clusterName string){
//	resp, _ := e.MyEtcdBackup.Client.Get(context.Background(), "/"+clusterName+"/", clientv3.WithPrefix())
//	key_num := len(resp.Kvs)
//
//	var wait sync.WaitGroup
//	wait.Add(key_num)
//
//	for _, ev := range resp.Kvs {
//		go func(key, val string){
//			last_index := len(clusterName) + 1
//			key = key[last_index:]
//			defer wait.Done()
//			_, error := e.MyEtcdMap[clusterName].Client.Put(context.TODO(), key, val)
//			if error != nil{
//				for error != nil {
//					fmt.Println(error)
//					_, error = e.MyEtcdMap[clusterName].Client.Put(context.TODO(), key, val)
//				}
//			}
//			fmt.Println("RESTORE["+clusterName+"] "+ " : ", key)
//		}(string(ev.Key), string(ev.Value))
//	}
//	wait.Wait()
//}

