package util

import (
	"fmt"
	"io/ioutil"
	"openmcp/openmcp/omcplog"
	"path"
	"strings"
)

//GetTemplate ./template/FILENAME 의 경로에 있는 text 를 가져오는 함수.
func GetTemplate(fileName string) (string, error) {
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
	filepath := "/root/template" //컨테이너의 해당 경로에 존재.
	shellPath := path.Join(filepath, fileName)

	file, err := ioutil.ReadFile(shellPath)
	if err != nil {
		omcplog.Error("Not ReadFile: "+shellPath, err)
		return "", err
	}
	if len(string(file)) == 0 {
		omcplog.Error("template is 0 size  : " + shellPath)
		return "", fmt.Errorf("template is 0 size  : " + shellPath)
	}
	return string(file), nil
}

func GetSnapshotTemplate(snapshotTime string, mountPath string) (string, error) {
	omcplog.V(5).Info("--- GetSnapshotTemplate start")
	cmd, err := GetTemplate("volumesnapshot.sh")
	if err != nil {
		omcplog.Error("get Template error : ", err)
		return "", err
	}

	ret1 := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	ret := strings.ReplaceAll(ret1, "!PATH", mountPath)

	omcplog.V(5).Info("--- GetSnapshotTemplate end")
	return ret, nil
}

func GetSnapshotRestoreTemplate(snapshotTime string, mountPath string) (string, error) {
	omcplog.V(5).Info("--- GetSnapshotRestoreTemplate start")
	cmd, err := GetTemplate("volumesnapshot-restore.sh")
	if err != nil {
		omcplog.Error("get Template error : ", err)
		return "", err
	}

	ret1 := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	ret := strings.ReplaceAll(ret1, "!PATH", mountPath)

	omcplog.V(5).Info("--- GetSnapshotRestoreTemplate end")
	return ret, nil
}

func GetLoopForSuccessTemplate() (string, error) {
	omcplog.V(5).Info("--- GetLoopForSuccessTemplate start")
	cmd, err := GetTemplate("loopForSuccess.sh")
	if err != nil {
		omcplog.Error("get Template error : ", err)
		return "", err
	}

	omcplog.V(5).Info("--- GetLoopForSuccessTemplate end")
	return cmd, nil
}
