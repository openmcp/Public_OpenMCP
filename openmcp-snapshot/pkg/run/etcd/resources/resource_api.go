package resources

import (
	"k8s.io/client-go/kubernetes"
)

const (
	// SnapshotTailName 스냅샷으로 생성될 서비스의 뒤에 붙는 이름
	SnapshotTailName = "-snapshot"
)

// Resource : 리소스 공통 인터페이스
type Resource interface {
	CreateResourceForJSON(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error)
	GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error)
}
