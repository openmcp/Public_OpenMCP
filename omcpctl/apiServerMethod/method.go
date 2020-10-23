package apiServerMethod

import (
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strings"
	"time"
)

var APP_KEY = "openmcp-apiserver"

func saveTokenToFile(token string) {
	filename := "/var/lib/omcpctl/token"
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = f.Truncate(0)
	_, err = fmt.Fprintln(f, "token: " + token)
	if err != nil {
		fmt.Println(err)
	}
}
func getTokenWithFile() (string, error){
	// token file
	filename := "/var/lib/omcpctl/token"

	// If Token File is not Exist
	// Then Get New Token
	_, err := os.Stat(filename)
	if os.IsNotExist(err){
		fmt.Println(filename+ " file is not")
		return "", cobrautil.NewError("File not exist")
	}

	tokenMap := make(map[string]string)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tokenMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	// If File Exist But Token Not Exist
	//Then Get New Token
	if _, ok := tokenMap["token"]; !ok {
		return "", cobrautil.NewError("token not exist")
	}

	// Validate the validity of the token.
	tokenString := tokenMap["token"]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(APP_KEY), nil
	})

	// ... error handling
	if err != nil {
		return "", err
	}
	_ = token

	// do something with decoded claims
	// exp is a value that indicates the validity period of the token
	exp := time.Duration(claims["exp"].(float64)) * time.Nanosecond

	exp_sec := int(exp.Seconds()*1000000000)
	cur_sec := int(time.Now().Unix())


	// If the expiration time is 300 seconds before, a new token is issued.
	if exp_sec - cur_sec < 300 {
		return "", cobrautil.NewError("Token is expired")
	}

	return tokenString, nil
}
func getToken() (string, error) {
	// Access to previously issued tokens
	// Check the expiration time and use the token if it is not expired
	// New tokens issued 5 minutes before expiration
	token, err := getTokenWithFile()
	if err == nil {
		return token, nil
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	LINK := "http://"+cobrautil.OpenMCPAPIServer+"/token"

	data := url.Values{}
	data.Set("username", "openmcp")
	data.Set("password", "keti")


	req, err := http.NewRequest("POST", os.ExpandEnv(LINK), strings.NewReader(data.Encode()))
	if err != nil {
		// handle err
		log.Fatalln(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Fatalln(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var prettyYaml map[string]string
	err = yaml.Unmarshal(body, &prettyYaml)
	if err != nil {
		panic(err.Error())
		return "", err
	}

	if _, ok := prettyYaml["token"]; !ok{
		return "", cobrautil.NewError("Cannot Get Token")
	}

	saveTokenToFile(prettyYaml["token"])
	return prettyYaml["token"], nil



}
func GetAPIServer(LINK string) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", os.ExpandEnv(LINK), nil)
	if err != nil {
		// handle err
		log.Fatalln(err)
	}
	TOKEN, err := getToken()

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization","Bearer "+ TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Check3", err)
		panic(err.Error())
	}
	return body, nil
}
func DeleteAPIServer(LINK string, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", LINK, body)
	if err != nil {
		// handle err
		log.Fatalln(err)
	}
	TOKEN, err := getToken()

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/yaml")
	req.Header.Set("Authorization", "Bearer "+ TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Check3", err)
		panic(err.Error())
	}

	return _body, nil

}
func PostAPIServer(LINK string, body io.Reader) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}


	req, err := http.NewRequest("POST", os.ExpandEnv(LINK), body)
	if err != nil {
		return nil, err
	}
	TOKEN, err := getToken()

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/yaml")
	req.Header.Set("Authorization", "Bearer " + TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Check3", err)
		panic(err.Error())
	}

	return _body, nil
}
func PutAPIServer(LINK string, body io.Reader) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}


	req, err := http.NewRequest("PUT", os.ExpandEnv(LINK), body)
	if err != nil {
		// handle err
		log.Fatalln(err)
	}
	TOKEN, err := getToken()

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/yaml")
	req.Header.Set("Authorization", "Bearer "+TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Check3", err)
		panic(err.Error())
	}

	return _body, nil
}