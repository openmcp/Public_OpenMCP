package registry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"openmcp/openmcp/omcplog"
	"testing"
)

func (r *RegistryManager) Test_ListGlobalRegistryImageTag(t *testing.T) ([]string, error) {
	imageName := "asdasd"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	getImageTagListURL := r.getImageTagListURL(imageName)
	omcplog.V(3).Info(getImageTagListURL)
	//r.getList(ListURL)
	req, requsetErr := http.NewRequest(http.MethodGet, getImageTagListURL, nil)
	if requsetErr != nil {
		// handle error
		return nil, requsetErr
	}
	//req.Header.Set("X-Registry-Auth", r.getAuthKey())
	resp, doErr := client.Do(req)
	if doErr != nil {
		// handle error
		return nil, doErr
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status + " error")
	}

	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(respBody)
	omcplog.V(3).Info(str)

	//Repositories output Structure
	type RepositoryInfo struct {
		Name string
		Tags []string
	}
	var output RepositoryInfo
	err = json.Unmarshal([]byte(str), &output)
	if err != nil {
		return nil, err
	}

	for _, tag := range output.Tags {
		//TagListURL := r.getImageTagListURL(repository)
		omcplog.V(3).Info(tag)
	}

	return output.Tags, nil
}

func Test_atestaaa() {

}
