package util

/*User Config*/
//external NFS IP
const EXTERNAL_NFS = "211.45.109.210"
const EXTERNAL_NFS_PATH = "/home/nfs/pv"
const EXTERNAL_NFS_NAME_PVC = "nfs-pvc"
const EXTERNAL_NFS_NAME_PV = "nfs-pv"

/*System Config*/
/***********************************************************/

//file copy cmd
const COPY_CMD = "cp -r"
const MKDIR_CMD = "mkdir -p"

//resource Type
const PVC = "PersistentVolumeClaim"
const PV = "PersistentVolume"
const DEPLOY = "Deployment"
const SERVICE = "Service"

//openmcp Namespace
const NameSpace = "openmcp"
