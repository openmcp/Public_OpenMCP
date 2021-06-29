package utils

// GlobalRepo 는 Global Repository 접속정보 입니다.
// User Config
var DockerHubRepo = globalRepoParamsMap{
	URI: "",
	//URI:      "index.docker.io",
	//LoginURL: "https://index.docker.io",
	LoginURL: "",
	Username: "openmcp",
	Password: "!Indy5515",

	//Cert: "-----BEGIN CERTIFICATE-----\nMIIDIDCCAgigAwIBAgIJALVEk3TZmsMlMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\nBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX\naWRnaXRzIFB0eSBMdGQwHhcNMjAwNDI0MDcyMzA4WhcNMjIwNDI0MDcyMzA4WjBF\nMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50\nZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAqN9zVhVQU+vMmO1UcXprZQmpcXTq+KFD05uyJPs4x7fiKglOE3gT2oXT\n7yuCNFTeHChybRB3E8/AOQSnHKJTpDy1WeGiGZ1DJKwSpzpRWSRGF7hjQlKSLNG+\n9VCrnIyJ6f7975yv1+1lUhgSq0c60JjkFkHNaNKO0ny1gKxAcWCNwZFPBojnx7r6\ncLXTWblROJduOmnfRZZHle1kQyXJKNBu1CXQXbGEJeN/YiMdE1nueeITbt7mTH79\nyXQZY2F2q265aKag3KReIG5FmcWqZ47FD/hMIttsN9u9btia6Vup2jwXCcWb8i99\nctfbIpMjX81fGRl6tbXDWl7f11e4SwIDAQABoxMwETAPBgNVHREECDAGhwQKAADg\nMA0GCSqGSIb3DQEBCwUAA4IBAQBJbzlRub6KMy79eyZb2eBYIxyf0KtLVA+LTlZe\nBqYfCS7yZ31XPOtzxpw251ji1A+k00Y/tO7oFF+ixEAPUT7XqijrcbM82qlzo5on\nVxf6ddfffppmXTa3miP8W69y7kEeBGtp3U6HIwsQbA4WOdttFvNf3LK7TswiuPYK\np7av7VLUEt8ndzydzAIQNhtalgb0oVVDRr6GahU/V0wHbIzlAwZS4byKyf117BPk\nIHtmp3WHockyobQVvYzqtbYvQ5+v55OUliPUIoKMZGIelaBQFdgeDgztOWPmQ3g9\nfZqlkN6xlj942sVHJfg4cV9GWR0+npgXMs/6sBFppN1XdOCW\n-----END CERTIFICATE-----",
}

var GlobalRepo = globalRepoParamsMap{
	URI:      "10.0.0.224:4999",
	LoginURL: "https://nanum:nanumrltnf@10.0.0.224:4999",
	Username: "nanum",
	Password: "nanumrltnf",
	//RegistryAuthKey: "ewogICJ1c2VybmFtZSI6ICJuYW51bSIsCiAgInBhc3N3b3JkIjogIm5hbnVtcmx0bmYiLAogICJlbWFpbCI6ICJkZXZAbmFudW0uY28ua3IiLAogICJzZXJ2ZXJhZGRyZXNzIjogIjEwLjAuMC4yMjQ6NTAwMCIKfQ==",

	Cert: "-----BEGIN CERTIFICATE-----\nMIIDIDCCAgigAwIBAgIJALVEk3TZmsMlMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\nBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX\naWRnaXRzIFB0eSBMdGQwHhcNMjAwNDI0MDcyMzA4WhcNMjIwNDI0MDcyMzA4WjBF\nMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50\nZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEAqN9zVhVQU+vMmO1UcXprZQmpcXTq+KFD05uyJPs4x7fiKglOE3gT2oXT\n7yuCNFTeHChybRB3E8/AOQSnHKJTpDy1WeGiGZ1DJKwSpzpRWSRGF7hjQlKSLNG+\n9VCrnIyJ6f7975yv1+1lUhgSq0c60JjkFkHNaNKO0ny1gKxAcWCNwZFPBojnx7r6\ncLXTWblROJduOmnfRZZHle1kQyXJKNBu1CXQXbGEJeN/YiMdE1nueeITbt7mTH79\nyXQZY2F2q265aKag3KReIG5FmcWqZ47FD/hMIttsN9u9btia6Vup2jwXCcWb8i99\nctfbIpMjX81fGRl6tbXDWl7f11e4SwIDAQABoxMwETAPBgNVHREECDAGhwQKAADg\nMA0GCSqGSIb3DQEBCwUAA4IBAQBJbzlRub6KMy79eyZb2eBYIxyf0KtLVA+LTlZe\nBqYfCS7yZ31XPOtzxpw251ji1A+k00Y/tO7oFF+ixEAPUT7XqijrcbM82qlzo5on\nVxf6ddfffppmXTa3miP8W69y7kEeBGtp3U6HIwsQbA4WOdttFvNf3LK7TswiuPYK\np7av7VLUEt8ndzydzAIQNhtalgb0oVVDRr6GahU/V0wHbIzlAwZS4byKyf117BPk\nIHtmp3WHockyobQVvYzqtbYvQ5+v55OUliPUIoKMZGIelaBQFdgeDgztOWPmQ3g9\nfZqlkN6xlj942sVHJfg4cV9GWR0+npgXMs/6sBFppN1XdOCW\n-----END CERTIFICATE-----",
}

// System Config
const (
	//ProjectDomain 프로젝트 도메인  - label 적용.
	ProjectDomain = "openmcp"

	//ProjectNamespace 프로젝트 네임스페이스 - kubernetes 에서 namespace 를 이것으로한다.
	ProjectNamespace = "openmcp"
)

type globalRepoParamsMap struct {
	URI      string
	LoginURL string
	Username string
	Password string
	//RegistryAuthKey string
	Cert string
}
