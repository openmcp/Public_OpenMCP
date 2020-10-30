
// convertResourceObj : json String 을 obj 로 변환
func convertResourceObj(resourceInfoJSON string) (*unstructured.Unstructured, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var unstructured *unstructured.Unstructured
	jsonEerr := json.Unmarshal(jsonBytes, &unstructured)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return unstructured, nil
}
