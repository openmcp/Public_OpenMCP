package util

import (
	"io/ioutil"
	"openmcp/openmcp/omcplog"
	"path"
	"strings"
)

//GetTemplate ./template/FILENAME 의 경로에 있는 text 를 가져오는 함수.
func GetTemplate(fileName string) string {
	omcplog.V(5).Info("--- GetSnapshotTemplate start")
	// _, filepath, _, ok := runtime.Caller(0)
	// if !ok {
	// 	panic("No caller information")
	// }
	// filepath, err := os.Getwd()
	// if err != nil {
	// 	omcplog.V(4).Info("Get File Path Error!")
	// }
	// filepath, err := os.Executable()
	// if err != nil {
	// 	omcplog.V(4).Info("Get File Path Error!")
	// }

	//projectPath := strings.Split(path.Dir(filepath), "/")
	//pathLen := len(projectPath) - 2
	filepath := "/root/template"
	shellPath := path.Join(filepath, fileName)

	file, err := ioutil.ReadFile(shellPath)
	if err != nil {
		panic("Not ReadFile: " + shellPath)
	}
	omcplog.V(5).Info("--- GetSnapshotTemplate end")
	return string(file)
}
func GetSnapshotTemplate(snapshotTime string, mountPath string) string {
	cmd := GetTemplate("volumesnapshot.sh")

	ret1 := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	ret := strings.ReplaceAll(ret1, "!PATH", mountPath)

	return ret
}
func GetSnapshotRestoreTemplate(snapshotTime string, mountPath string) string {
	cmd := GetTemplate("volumesnapshot-restore.sh")

	ret1 := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	ret := strings.ReplaceAll(ret1, "!PATH", mountPath)

	return ret
}
