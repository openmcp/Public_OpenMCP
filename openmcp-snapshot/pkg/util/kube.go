package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"openmcp/openmcp/omcplog"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Event struct {
	Name          string
	LastTimestamp metav1.Time
	Reason        string
	Message       string
}

//GetClientset kubernetes의 clientset 생성
func GetClientset(clusterInfo string) (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		return nil, err
	}
	clientconf, err := con.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err = kubernetes.NewForConfig(clientconf)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Event 리소스에서 해당 리소스를 찾아 SuccessfulCreate이 아닐 경우 해당 상세 메시지를 Json 배열화 시키는 함수
func FindErrorForEvent(targetListClient *kubernetes.Clientset, resourceName string) (string, error) {
	omcplog.V(3).Info("-- error collect Event Resource")

	var regexpErr error
	matchResult := false
	errDetail := ""
	//시간초과 - 오류 루틴으로 진입

	//이벤트에서 해당 오류 찾아서 도출. Reason에 SuccessfulCreate가 포함된 경우는 CMD 에서 오류를 찾아야한다.
	events, err := targetListClient.CoreV1().Events(JOB_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	tmpEvents := []Event{}

	isSuccessfulCreate := false
	for _, event := range events.Items {
		matchResult, regexpErr = regexp.MatchString(resourceName, event.Name+"-") //job/sns-1629180186-1629180186-1-job가 아닌 pod/sns-1629180186-1629180186-1-job-9954d를 검출하기 위함
		if regexpErr == nil && matchResult == true {
			omcplog.V(3).Info(event.Reason)
			if event.Reason == "SuccessfulCreate" {
				isSuccessfulCreate = true
				omcplog.V(3).Info("-- job- running success")
				break
			}
			tmpEvent := Event{}
			tmpEvent.Name = event.Name
			tmpEvent.LastTimestamp = event.LastTimestamp
			tmpEvent.Message = event.Message
			tmpEvent.Reason = event.Reason
			tmpEvents = append(tmpEvents, tmpEvent)
		}
	}
	if !isSuccessfulCreate {
		jsonTmp, err := json.Marshal(tmpEvents)
		if err != nil {
			omcplog.V(3).Info(err, "-----------") //어쩔수 없음.
		}
		errDetail = string(jsonTmp)
	}
	return errDetail, nil
}

// pod 내에서 쉘 실행시키는 함수
func RunCommand(client kubernetes.Interface, config *restclient.Config, podName string, command string, namespace string) (*ExecutorResult, error) {
	omcplog.V(3).Info("----- Start Command -----")
	//omcplog.V(3).Info(command)

	cmd := []string{
		"bash",
		"-c",
		command,
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
	option := &corev1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	omcplog.V(3).Info("1. setConfig")
	//omcplog.V(3).Info(config)
	result := new(ExecutorResult)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		omcplog.Error("NewSPDYExecutor err : ", err)
		return result, err
	}

	//err = exec.Stream(remotecommand.StreamOptions{
	//	Stdin:  os.Stdin,
	//	Stdout: os.Stdout,
	//	Stderr: os.Stderr,
	//})
	//stdin := strings.NewReader("Hello, Reader!")
	//var stdout io.Writer
	//var stderr io.Writer

	omcplog.V(3).Info("2. Run stream")
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  &result.Stdin,
		Stdout: &result.Stdout,
		Stderr: &result.Stderr,
		//Stdin:  os.Stdin,
		//Stdout: os.Stdout,
		//Stderr: os.Stderr,
	})
	omcplog.V(3).Info("3. Print Error")

	omcplog.V(3).Info("3. Print Error")
	if err != nil {
		omcplog.Error(err)
		if strings.Contains(err.Error(), "100") {
			omcplog.Error("Command Error : TargetFile is empty!")
			return result, fmt.Errorf("Command Error : TargetFile is empty!")
		} else {
			omcplog.Error("NewSPDYExecutor Stream err : ", err)
			return result, err
		}
	}
	return result, nil
}

// ExecutorResult contains the outputs of the execution.
type ExecutorResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Stdin  bytes.Buffer
}
