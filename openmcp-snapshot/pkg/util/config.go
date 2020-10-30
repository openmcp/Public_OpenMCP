package util

/*User Config*/
//external NFS IP
const EXTERNAL_NFS = "211.45.109.210"
const EXTERNAL_NFS_PATH_STORAGE = "/home/nfs/storage"

//OpenMCP Master 아이디
const MASTER_IP = "10.0.0.226"

//external ETCD IP
const EXTERNAL_ETCD = "10.0.0.221:12379"

/*System Config*/
/***********************************************************/

//resource Type
const PVC = "PersistentVolumeClaim"
const PV = "PersistentVolume"
const DEPLOY = "Deployment"
const SERVICE = "Service"

//pv spec
const DEFAULT_VOLUME_SIZE = "10Gi"
