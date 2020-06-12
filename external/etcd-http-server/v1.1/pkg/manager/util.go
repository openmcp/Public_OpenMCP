package manager

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)
var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second

	SNAPSHOT_DIR = "/root/data/etcd/snapshot"
	RESTORE_DATA_DIR = "/root/data/etcd/data"
	TMP_DIR = "/root/data/etcd/tmp"
	//PKI_DIR ="/root/workspace/go/gopath/src/etcd-syncer/pki"
)

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}
func ensureDir(dirName string) error {

	err := os.MkdirAll(dirName, os.ModeDir)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
func FileNotExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true
	}
	return false
}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}

func LoadYAMLFile(file string) (*KubeConfig, error) {
	config := &KubeConfig{}
	yamlfile, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlfile, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
// TODO: TLS
func newEmbedURLs(n int) (urls []url.URL) {
	urls = make([]url.URL, n)
	for i := 0; i < n; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		u, _ := url.Parse(fmt.Sprintf("unix://localhost:%d", rand.Intn(45000)))
		urls[i] = *u
	}
	return urls
}
func CopyDir(src_path, dest_path string) error{
	cmdStr := "cp -r " + src_path + " " + dest_path
	_, err := CmdExec(cmdStr)
	if err != nil {
		fmt.Println("CopyErr !", err)
		return err
	}
	return nil

}
func MoveDir(src_path, dest_path string) error{
	cmdStr := "mv " + src_path + " " + dest_path
	_, err := CmdExec(cmdStr)
	if err != nil {
		fmt.Println("MoveErr !", err)
		return err
	}
	return nil

}
func DeleteDir(src_path string) {
	cmdStr := "rm -r " + src_path
	_, err := CmdExec(cmdStr)
	if err != nil {
		fmt.Println("DeleteErr !", err)
	}
}
func RemoteCopyDir(host, srcpath, destpath string){
	cmdStr := "scp -r "+srcpath +" root@"+host+":"+destpath
	_, err := CmdExec(cmdStr)
	if err != nil {
		fmt.Println("Err !", err)
	}
}
func RemoteDeleteDir(host, srcpath string){
	cmdStr := "ssh "+host+" rm -rf "+ srcpath
	_, err := CmdExec(cmdStr)
	if err != nil {
		fmt.Println("Err !", err)
	}
}

func CmdExec(cmdStr string) (string, error) {
	fmt.Println(cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func MapsKeyCompare(mapA, mapB map[string]string) ([]string, []string, []string){
	both_keys := []string{}
	onlyA_keys := []string{}
	onlyB_keys := []string{}

	for k, _ := range mapA{
		if _, ok := mapB[k]; ok {
			both_keys = append(both_keys, k)
		} else {
			onlyA_keys = append(onlyA_keys, k)
		}
	}
	for k, _ := range mapB{
		exist := false
		for _, both_key := range both_keys {
			if both_key == k{
				exist = true
				break
			}
		}
		if !exist {
			onlyB_keys = append(onlyB_keys, k)
		}
	}

	return both_keys, onlyA_keys, onlyB_keys

}
