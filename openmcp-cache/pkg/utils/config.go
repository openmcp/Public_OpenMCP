package utils

// System Config
const (
	//ProjectDomain 프로젝트 도메인  - label 적용.
	ProjectDomain = "openmcp"
	//ProjectNamespace 프로젝트 네임스페이스 - kubernetes 에서 namespace 를 이것으로한다.
	ProjectNamespace = "openmcp"
)

var ImageCacheNfs = "211.45.109.210"

type globalRepoParamsMap struct {
	URI      string
	LoginURL string
	Username string
	Password string
	//RegistryAuthKey string
	Cert string
}
