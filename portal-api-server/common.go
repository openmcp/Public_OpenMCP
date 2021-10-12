package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Resultmap struct {
	secs float64
	url  string
	data map[string]interface{}
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetOpenMCPToken() string {
	// caCert, err := ioutil.ReadFile(`server.crt`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			// TLSClientConfig: &tls.Config{
			// 	RootCAs: caCertPool,
			// },
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	account := Account{"openmcp", "keti"}

	pbytes, _ := json.Marshal(account)
	buff := bytes.NewBuffer(pbytes)
	// resp, err := client.Get("https://" + openmcpURL + "/token?username=openmcp&password=keti")
	resp, err := client.Post("https://"+openmcpURL+"/token", "application/json", buff)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var data map[string]interface{}
	token := ""
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(bodyBytes), &data)
		token = data["token"].(string)

	} else {
		fmt.Println("failed")
	}
	return token
}

func CallAPI(token string, url string, ch chan<- Resultmap) {
	start := time.Now()
	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	// var client http.Client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close() // 리소스 누출 방지
	if err != nil {
		// ch <- fmt.Sprintf("while reading %s: %v", url, err)
		// return
		log.Fatal(err)
	}
	json.Unmarshal([]byte(bodyBytes), &data)

	secs := time.Since(start).Seconds()

	// ch <- fmt.Sprintf("%.2fs %s %v", secs, url, data)

	ch <- Resultmap{secs, url, data}

}

func GetJsonBody(rbody io.Reader) map[string]interface{} {
	bodyBytes, err := ioutil.ReadAll(rbody)

	var data map[string]interface{}

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(bodyBytes), &data)
	return data
}

func PostYaml(url string, yaml io.Reader) ([]byte, error) {
	token := GetOpenMCPToken()
	// fmt.Println("yaml   :", yaml)
	var bearer = "Bearer " + token
	req, err := http.NewRequest("POST", url, yaml)

	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	// var client http.Client

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	str := string(respBody)
	fmt.Println(str)
	return respBody, nil

}

func GetStringElement(nMap interface{}, keys []string) string {
	result := ""

	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(string)
				} else {
					result = childMap.(string)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = "-"
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = "-"
					break
				}
			}
		}
	} else {
		result = "-"
	}
	return result
}

func GetIntElement(nMap interface{}, keys []string) int {
	result := 0
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(int)
				} else {
					result = childMap.(int)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = 0
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = 0
					break
				}
			}
		}
	} else {
		result = 0
	}
	return result
}

func GetFloat64Element(nMap interface{}, keys []string) float64 {
	var result float64 = 0.0
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0].(float64)
				} else {
					result = childMap.(float64)
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = 0.0
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = 0.0
					break
				}
			}
		}
	} else {
		result = 0.0
	}
	return result
}

func GetInterfaceElement(nMap interface{}, keys []string) interface{} {
	var result interface{}
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})[0]
				} else {
					result = childMap
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			}
		}
	} else {
		result = nil
	}
	return result
}

func GetArrayElement(nMap interface{}, keys []string) []interface{} {
	var result []interface{}
	if nMap.(map[string]interface{})[keys[0]] != nil {
		childMap := nMap.(map[string]interface{})[keys[0]]
		for i, _ := range keys {
			typeCheck := fmt.Sprintf("%T", childMap)

			if len(keys)-1 == i {
				if "[]interface {}" == typeCheck {
					result = childMap.([]interface{})
				} else {
					result = childMap.([]interface{})
				}
				break
			}

			if "[]interface {}" == typeCheck {
				if childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.([]interface{})[0].(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			} else {
				if childMap.(map[string]interface{})[keys[i+1]] != nil {
					childMap = childMap.(map[string]interface{})[keys[i+1]]
				} else {
					result = nil
					break
				}
			}
		}
	} else {
		result = nil
	}
	return result
}
