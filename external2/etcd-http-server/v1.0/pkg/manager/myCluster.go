package manager

type MyCluster struct {
	ClusterName        string
	IP                 string
	OpenMCPMasterIP    string
	isEtcdBackupServer bool
}
