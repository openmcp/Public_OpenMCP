package main

import (
	"etcd-syncer/pkg/manager"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type HttpManager struct {
	HTTPServer_IP   string
	HTTPServer_PORT string
	OpenMCP_IP      string
	stop            chan string
	backup_status   string
}

func (h *HttpManager) main(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Main")

	w.Write([]byte("------Help------\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/backup/start?time=5\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/backup/stop\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/backup/status\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/get/snapshot/list\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/restore/cluster?cluster=CLUSTER_NAME&file=SNAPSHOT.db\n"))
	w.Write([]byte("http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/restore/all?file=SNAPSHOT.db\n"))

}
func (h *HttpManager) etcd_backup_start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Backup Start")
	period, _ := strconv.Atoi(r.URL.Query().Get("time"))
	if period <= 0 {
		period = 5
	}

	if h.backup_status == "STOP" {
		fmt.Println("Etcd Backup Start")
		h.backup_status = "START"
		go func(w http.ResponseWriter, r *http.Request) {
			for {
				select {
				case <-h.stop:
					fmt.Println("STOP Recv")
					h.backup_status = "STOP"
					return
				default:
					fmt.Println("Wait " + strconv.Itoa(period) + " Seconds")
					time.Sleep(time.Duration(period) * time.Second)
					fmt.Println("Start Backup")
					etcdManager := manager.NewEtcdManager(h.OpenMCP_IP, h.HTTPServer_IP)
					etcdManager.MyFuncSyncAll()

					db_filename := fmt.Sprintf("snapshot%d.db", time.Now().Unix())
					for cluster_name, myEtcd := range etcdManager.MyEtcdMap {
						fmt.Println("[" + cluster_name + "] Backup Start")
						db_filefullname, err := myEtcd.SnapShot(db_filename)
						fmt.Println("[" + cluster_name + "] Backup End -> " + db_filefullname)
						if err != nil {
							fmt.Println(err)
							return
						}
					}

					fmt.Println("End Backup")
				}
			}
		}(w, r)
		w.Write([]byte("Etcd Backup Start"))
	} else {
		w.Write([]byte("Etcd Backup Already Started"))
	}

}
func (h *HttpManager) etcd_backup_stop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Backup Stop")
	if h.backup_status == "START" {
		h.stop <- "STOP"
		w.Write([]byte("Etcd Backup Stop"))
	} else {
		w.Write([]byte("Etcd Backup Already Stopped"))
	}
}

func (h *HttpManager) etcd_backup_status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Backup Status")
	w.Write([]byte(h.backup_status))
}

func (h *HttpManager) etcd_get_snapshot_list(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Get Snapshot List")

	snapshot_map := make(map[string][]string)

	dirs, err := ioutil.ReadDir(manager.SNAPSHOT_DIR)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range dirs {
		dirname := d.Name()

		files, err := ioutil.ReadDir(filepath.Join(manager.SNAPSHOT_DIR, dirname))
		if err != nil {
			fmt.Println(dirname)
			fmt.Println(err)
		}

		for _, f := range files {
			filename := f.Name()

			if filename[len(filename)-3:] != ".db" {
				continue
			}
			_, ok := snapshot_map[filename]
			if !ok {
				snapshot_map[filename] = []string{}
			}
			snapshot_map[filename] = append(snapshot_map[filename], dirname)

		}
	}

	keys := make([]string, 0, len(snapshot_map))
	for k := range snapshot_map {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		filename := k
		cluster_list := snapshot_map[k]

		unix_time, _ := strconv.Atoi(filename[8:18])
		rfc_time := time.Unix(int64(unix_time), 0).Format(time.RFC3339)

		w.Write([]byte(filename + " (" + rfc_time + ") ["))
		for i, cluster_name := range cluster_list {
			if i == len(cluster_list)-1 {
				w.Write([]byte(cluster_name))
			} else {
				w.Write([]byte(cluster_name + " / "))
			}
		}
		w.Write([]byte("]\n"))
	}
}
func (h *HttpManager) etcd_restore_all(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Restore All")

	if h.backup_status != "STOP" {
		w.Write([]byte("First, Must be Etcd Backup Server Stopped"))
		return
	}
	db_filename := r.URL.Query().Get("file")

	etcdManager := manager.NewEtcdManager(h.OpenMCP_IP, h.HTTPServer_IP)

	fileExist := false
	for _, myEtcd := range etcdManager.MyEtcdList {
		clusterName := myEtcd.MyCluster.ClusterName
		db_filepath := filepath.Join(manager.SNAPSHOT_DIR, clusterName, db_filename)
		if manager.FileExists(db_filepath) {
			fileExist = true
			break
		}
	}

	if !fileExist {
		w.Write([]byte("Argument file " + db_filename + " is not exist\n\n"))
		w.Write([]byte("First, you get list 'SNAPSHOT.db'\n"))
		w.Write([]byte("=> http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/get/snapshot/list\n"))
		w.Write([]byte("Second, you retry restore\n"))
		w.Write([]byte("=> http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/restore?file=SNAPSHOT.db\n"))
		return
	}
	for _, myEtcd := range etcdManager.MyEtcdList {
		clusterName := myEtcd.MyCluster.ClusterName
		db_filepath := filepath.Join(manager.SNAPSHOT_DIR, clusterName, db_filename)
		if manager.FileExists(db_filepath) {
			fmt.Println("[" + clusterName + "] Restore Start")
			db_datadir := etcdManager.MyEtcdMap[clusterName].Restore(db_filepath)
			manager.RemoteDeleteDir(etcdManager.MyEtcdMap[clusterName].MyCluster.IP, "/var/lib/etcd")
			manager.RemoteCopyDir(etcdManager.MyEtcdMap[clusterName].MyCluster.IP, db_datadir, "/var/lib/etcd")

		}

	}

	w.Write([]byte("Done"))

}
func (h *HttpManager) etcd_restore_cluster(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect Etcd Restore Cluster")

	if h.backup_status != "STOP" {
		w.Write([]byte("First, Must be Etcd Backup Server Stopped"))
		return
	}
	db_filename := r.URL.Query().Get("file")
	clusterName := r.URL.Query().Get("cluster")

	db_filepath := filepath.Join(manager.SNAPSHOT_DIR, clusterName, db_filename)

	if manager.FileNotExists(db_filepath) {
		w.Write([]byte("Argument file " + db_filename + " is not exist\n\n"))
		w.Write([]byte("First, you get list 'SNAPSHOT.db'\n"))
		w.Write([]byte("=> http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/get/snapshot/list\n"))
		w.Write([]byte("Second, you retry restore\n"))
		w.Write([]byte("=> http://" + h.HTTPServer_IP + ":" + h.HTTPServer_PORT + "/etcd/restore?file=SNAPSHOT.db\n"))
		return
	}

	etcdManager := manager.NewEtcdManager(h.OpenMCP_IP, h.HTTPServer_IP)

	db_datadir := etcdManager.MyEtcdMap[clusterName].Restore(db_filepath)

	fmt.Println(db_datadir)

	manager.RemoteDeleteDir(etcdManager.MyEtcdMap[clusterName].MyCluster.IP, "/var/lib/etcd")
	manager.RemoteCopyDir(etcdManager.MyEtcdMap[clusterName].MyCluster.IP, db_datadir, "/var/lib/etcd")

	w.Write([]byte("Done"))

}
func (h *HttpManager) test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect test")
	//filename := "test.db"
	db_filename := fmt.Sprintf("snapshot%d.db", time.Now().Unix())

	etcdManager := manager.NewEtcdManager(h.OpenMCP_IP, h.HTTPServer_IP)
	db_filename, err := etcdManager.MyEtcdMap["cluster3"].SnapShot(db_filename)
	if err != nil {
		fmt.Println(nil)
		return
	}
	//etcdManager.MyEtcdMap["cluster3"].DelWithPrefix("/")
	//time.Sleep(30 * time.Second)

	etcdManager.MyEtcdMap["cluster3"].Restore(db_filename)
}

func main() {
	HTTPServer_IP := "10.0.3.12"
	HTTPServer_PORT := "8090"
	OpenMCP_IP := "10.0.3.20"

	httpManager := &HttpManager{
		HTTPServer_IP:   HTTPServer_IP,
		HTTPServer_PORT: HTTPServer_PORT,
		OpenMCP_IP:      OpenMCP_IP,
		stop:            make(chan string, 1),
		backup_status:   "STOP",
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/", httpManager.main)
	handler.HandleFunc("/etcd/backup/start", httpManager.etcd_backup_start)
	handler.HandleFunc("/etcd/backup/stop", httpManager.etcd_backup_stop)
	handler.HandleFunc("/etcd/backup/status", httpManager.etcd_backup_status)
	handler.HandleFunc("/etcd/get/snapshot/list", httpManager.etcd_get_snapshot_list)
	handler.HandleFunc("/etcd/restore/cluster", httpManager.etcd_restore_cluster)
	handler.HandleFunc("/etcd/restore/all", httpManager.etcd_restore_all)
	handler.HandleFunc("/test", httpManager.test)

	server := &http.Server{Addr: HTTPServer_IP + ":" + HTTPServer_PORT, Handler: handler}
	server.ListenAndServe()

	//for {
	//	fmt.Println("Start Backup")
	//	etcdManager := NewEtcdManager()
	//etcdManager.MyFuncBackAll()
	//etcdManager.MyFuncSnapShot()
	//	fmt.Println("End Backup")
	//	time.Sleep(5 * time.Second)
	//
	//}

	//etcdManager.MyFuncWatch()
	//etcdManager.MyFuncBackup("")
	//etcdManager.MyFuncRestore("", "")
	//etcdManager.MyFuncGet("")
	//etcdManager.MyFuncPut("","")
	//etcdManager.MyFuncDel("")

	//etcdManager.MyFuncBackup("/registry/csinodes/kube3-worker2")
	//etcdManager.MyFuncBackup("/registry/leases/kube-node-lease/kube3-worker2")
	//etcdManager.MyFuncBackup("/registry/minions/kube3-worker2")
	//etcdManager.MyFuncBackAll()

}
