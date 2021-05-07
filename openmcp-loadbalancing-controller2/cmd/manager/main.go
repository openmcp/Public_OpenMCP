/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"sync"
	//"time"

	"bytes"
	"fmt"
	"io/ioutil"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	"net/http"
)


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		//w.Write([]byte("hello world\n"))
		// Request 객체 생성
		//fmt.Println(req)
		proxyScheme := "http"
		//proxyHost := "localhost:9091"
		proxyHost := "10.0.3.196"

		// we need to buffer the body if we want to read it here and send it
		// in the request.
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// you can reassign the body if you need to parse it as multipart
		req.Body = ioutil.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		url := fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, req.RequestURI)

		proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
		//필요시 헤더 추가 가능
		proxyReq.Header = make(http.Header)
		for h, val := range req.Header {
			proxyReq.Header[h] = val
		}
		proxyReq.RemoteAddr = req.RemoteAddr
		//fmt.Println(proxyReq)

		// Client객체에서 Request 실행
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// 결과 출력
		bytes, _ := ioutil.ReadAll(resp.Body)
		str := string(bytes) //바이트를 문자열로
		//fmt.Println(str)
		w.Write([]byte(str))


	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
