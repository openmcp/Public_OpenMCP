package util

import (
	"encoding/json"
	"fmt"
)

// Obj2JsonString : Deployment 등과 같은 interface 를 json string 으로 변환.
func Obj2JsonString(obj interface{}) (string, error) {

	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println("===Obj2JsonString===")
	fmt.Println(string(json))

	return string(json), nil
}
