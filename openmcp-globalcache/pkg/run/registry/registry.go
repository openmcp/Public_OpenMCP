package registry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"openmcp/openmcp/omcplog"
)

//https://github.com/docker/go-docker/blob/master/ 참조하여 작성

//ListGlobalRegistryImage 글로벌 레지스트리의 이미지 출력.
func (r *RegistryManager) ListGlobalRegistryImage() ([]string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	getListURL := r.getListURL()
	omcplog.V(3).Info(getListURL)
	//r.getList(ListURL)
	req, requsetErr := http.NewRequest(http.MethodGet, getListURL, nil)
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
	type Repositories struct {
		Repositories []string
	}
	var output Repositories
	err = json.Unmarshal([]byte(str), &output)
	if err != nil {
		return nil, err
	}

	for _, repository := range output.Repositories {
		//TagListURL := r.getImageTagListURL(repository)
		omcplog.V(3).Info(repository)
	}

	return output.Repositories, nil
}

//ListGlobalRegistryImageTag 글로벌 레지스트리의 이미지에 대한 태그 목록 출력
func (r *RegistryManager) ListGlobalRegistryImageTag(imageName string) ([]string, error) {

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

//DeleteGlobalRegistryImage 는 글로벌 레지스트리의 이미지 삭제 (해당 태그)
func (r *RegistryManager) DeleteGlobalRegistryImage(repository string, tag string) (bool, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	getHeaderURL := r.getRequestDeleteHeaderURL(repository, tag)
	omcplog.V(3).Info(getHeaderURL)
	req, requsetErr := http.NewRequest(http.MethodGet, getHeaderURL, nil)
	if requsetErr != nil {
		// handle error
		return false, requsetErr
	}
	//req.Header.Set("X-Registry-Auth", r.getAuthKey())
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	resp, doErr := client.Do(req)
	if doErr != nil {
		// handle error
		return false, doErr
	}
	if resp.StatusCode != http.StatusOK {

		defer resp.Body.Close()
		// Response 체크.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		str := string(respBody)
		omcplog.V(3).Info(str)

		return false, errors.New(resp.Status + ": " + str)
	}

	digest := resp.Header.Get("Docker-Content-Digest")
	authorization := resp.Request.Header.Get("Authorization")

	// ---

	getDeleteURL := r.getDeleteURL(repository, digest)
	omcplog.V(3).Info(getDeleteURL)
	req, requsetErr = http.NewRequest(http.MethodDelete, getDeleteURL, nil)
	if requsetErr != nil {
		// handle error
		return false, requsetErr
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, doErr = client.Do(req)
	if doErr != nil {
		// handle error
		return false, doErr
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {

		defer resp.Body.Close()
		// Response 체크.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		str := string(respBody)
		omcplog.V(3).Info(str)

		return false, errors.New(resp.Status + ": " + str)
	}

	return true, nil
}

//DeleteGlobalRegistryImageAllTag 모든 태그를 삭제.
func (r *RegistryManager) DeleteGlobalRegistryImageAllTag(repository string) (bool, error) {

	imageTags, err := r.ListGlobalRegistryImageTag(repository)
	if err != nil {
		// handle error
		return false, err
	}

	var isSuccess bool
	for _, imageTag := range imageTags {
		isSuccess, err = r.DeleteGlobalRegistryImage(repository, imageTag)
		if !isSuccess || err != nil {
			// handle error
			return false, err
		}

	}

	return isSuccess, nil
}
