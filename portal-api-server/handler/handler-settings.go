package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetOpenmcpPolicy(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	openmcpPolicyUrl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/openmcppolicys?clustername=openmcp"
	go CallAPI(token, openmcpPolicyUrl, ch)
	openmcpPolicy := <-ch
	openmcpPolicyData := openmcpPolicy.data

	openmcpPolicyRes := OpenmcpPolicyRes{}
	openmcpPolicyInfo := OpenmcpPolicy{}

	//get clusters Information
	for _, element := range openmcpPolicyData["items"].([]interface{}) {
		policyName := GetStringElement(element, []string{"metadata", "name"})
		policyStatus := GetStringElement(element, []string{"spec", "policyStatus"})

		value := ""
		policies := GetArrayElement(element, []string{"spec", "template", "spec", "policies"})
		if policies != nil {
			for _, item := range policies {
				policyType := GetStringElement(item, []string{"type"})
				policyValues := GetArrayElement(item, []string{"value"})
				if policyValues != nil {
					for _, item := range policyValues {
						policyValueStr := fmt.Sprintf("%v", item)
						value = value + policyType + " : " + policyValueStr + "|"
						// if j+1 == len(policies) && len(policyValues) == 1 {
						// 	value = value + policyType + " : " + policyValueStr
						// } else if j+1 == len(policies) && len(policyValues) > 1 {
						// 	if j+1 == len(policyValues) {
						// 		value = value + policyType + " : " + policyValueStr

						// 	} else {
						// 		value = value + policyType + " : " + policyValueStr + "|"
						// 	}
						// } else {
						// 	value = value + policyType + " : " + policyValueStr + "|"
						// }
					}
				}
			}

			openmcpPolicyInfo.Name = policyName
			openmcpPolicyInfo.Status = policyStatus
			openmcpPolicyInfo.Value = value

			openmcpPolicyRes.OpenmcpPolicy = append(openmcpPolicyRes.OpenmcpPolicy, openmcpPolicyInfo)
		}
	}
	json.NewEncoder(w).Encode(openmcpPolicyRes.OpenmcpPolicy)
}

func UpdateOpenmcpPolicy(w http.ResponseWriter, r *http.Request) {
	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	editPolicyRes := EditPolicyRes{}

	policyName := data["policyName"].(string)
	policyValues := data["values"].([]interface{})

	for _, element := range policyValues {
		editInfo := EditInfo{}
		editInfo.Op = GetStringElement(element, []string{"op"})
		editInfo.Path = GetStringElement(element, []string{"path"})
		editInfo.Value = GetStringElement(element, []string{"value"})

		editPolicyRes.EditPolicy = append(editPolicyRes.EditPolicy, editInfo)
	}

	// var body []interface{}
	// body = append(body, EditInfo{"replace", "/spec/metalLBRange/addressFrom", ""})
	// body = append(body, EditInfo{"replace", "/spec/metalLBRange/addressTo", ""})
	// body = append(body, EditInfo{"replace", "/spec/joinStatus", "UNJOIN"})

	var jsonErrs []jsonErr

	// https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcppolicys/log-level?clustername=openmcp
	projectURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcppolicys/" + policyName + "?clustername=" + openmcpClusterName

	resp, err := CallPatchAPI(projectURL, "application/json-patch+json", editPolicyRes.EditPolicy, true)
	var msg jsonErr

	if err != nil {
		msg = jsonErr{503, "failed", "request fail"}
	}

	var dataRes map[string]interface{}
	json.Unmarshal([]byte(resp), &dataRes)
	if dataRes != nil {
		if dataRes["kind"].(string) == "Status" {
			msg = jsonErr{501, "failed", dataRes["message"].(string)}
		} else {
			msg = jsonErr{200, "success", "OpenmcpPolicy Update Completed"}
		}
	}

	jsonErrs = append(jsonErrs, msg)

	json.NewEncoder(w).Encode(jsonErrs)
}
