/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"
	"fmt"
	"log"
	"openmcp/openmcp/util"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func RunCache() error {
	// func RunCache(clusterIP string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.0.0.226:12379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := cli.Get(ctx, "cache", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	var m map[string]int
	m = make(map[string]int)
	for _, ev := range resp.Kvs {
		keys := strings.Split(strings.TrimSpace(string(ev.Key)), "/")
		key := ""
		for i := 2; i < len(keys); i++ {
			key += "/" + keys[i]
		}
		value, _ := strconv.Atoi(strings.TrimSpace(string(ev.Value)))
		m[key] += value
	}
	result := util.CmdExec2("kubectl get nodes -c 10.0.0.227")
	fmt.Println(result)
	// for key, val := range m {
	// 	if val >= 6 {
	// 		fmt.Println(key, val)
	// 	}
	// }
	return err
}
