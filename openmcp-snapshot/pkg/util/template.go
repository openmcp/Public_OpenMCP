package util

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
)

//GetTemplate ./template/FILENAME 의 경로에 있는 text 를 가져오는 함수.
func GetTemplate(fileName string) string {
	fmt.Printf("--- GetSnapshotTemplate start")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	shellPath := path.Join(path.Dir(filename), "template", fileName)
	fmt.Printf("Filename : %q, Dir : %q\n", filename, shellPath)
	file, err := ioutil.ReadFile(shellPath)
	if err != nil {
		panic("Not ReadFile: " + shellPath)
	}
	fmt.Printf("--- GetSnapshotTemplate end")
	return string(file)
}
func GetSnapshotTemplate(snapshotTime string) string {
	cmd := GetTemplate("volumesnapshot.sh")

	ret := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	return ret
}
func GetSnapshotRestoreTemplate(snapshotTime string) string {
	cmd := GetTemplate("volumesnapshot-restore.sh")

	ret := strings.ReplaceAll(cmd, "!DATE", snapshotTime)
	return ret
}
