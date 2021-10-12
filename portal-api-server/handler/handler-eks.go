package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func StopEKSNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// region := r.URL.Query().Get("region")
	// nodename := r.URL.Query().Get("node")
	// // nodegroup := r.URL.Query().Get("nodegroup")
	// // desiredSizeStr := r.URL.Query().Get("nodecount")
	// // http://192.168.0.51:4885/apis/eksinstancestop?region=ap-northeast-2&node=ip-172-31-58-160.ap-northeast-2.compute.internal
	// akid := "AKIAJGFO6OXHRN2H6DSA"                          //
	// secretkey := "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" //

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	region := data["region"].(string)
	nodename := data["node"].(string)
	akid := data["akid"].(string)
	secretkey := data["secretKey"].(string)

	// fmt.Println(region)
	// fmt.Println(nodename)
	// fmt.Println(akid)
	// fmt.Println(secretkey)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region), //
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}
	svc := ec2.New(sess)
	lists, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	var targetID string
	for _, v := range lists.Reservations {
		for _, e := range v.Instances {
			if *e.PrivateDnsName == nodename {
				targetID = *e.InstanceId
				break
			}
		}
	}

	asSvc := autoscaling.New(sess)
	var targetASG string
	instances, _ := asSvc.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{})
	for _, instance := range instances.AutoScalingInstances {
		if *instance.InstanceId == targetID {
			targetASG = *instance.AutoScalingGroupName
			break
		}
	}

	// fmt.Println("targetID", targetID, targetASG)
	input := &autoscaling.EnterStandbyInput{
		AutoScalingGroupName: aws.String(targetASG),
		InstanceIds: []*string{
			aws.String(targetID),
		},
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	}

	result, err := asSvc.EnterStandby(input)
	fmt.Println(result)

	if err != nil {
		errmsg := jsonErr{503, "failed", err.Error()}
		json.NewEncoder(w).Encode(errmsg)
	} else {
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(targetID),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StopInstances(input)
			if err != nil {
				errmsg := jsonErr{503, "failed", err.Error()}
				json.NewEncoder(w).Encode(errmsg)
			} else {
				fmt.Println("Success", result.StoppingInstances)
				errmsg := jsonErr{200, "success", "vm instance stopping"}
				json.NewEncoder(w).Encode(errmsg)
			}
		} else {
			errmsg := jsonErr{503, "failed", err.Error()}
			json.NewEncoder(w).Encode(errmsg)
		}
	}
}

func StartEKSNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	region := data["region"].(string)
	nodename := data["node"].(string)
	akid := data["akid"].(string)
	secretkey := data["secretKey"].(string)

	// fmt.Println(region)
	// fmt.Println(nodename)
	// fmt.Println(akid)
	// fmt.Println(secretkey)

	// region := r.URL.Query().Get("region")
	// nodename := r.URL.Query().Get("node")
	// // nodegroup := r.URL.Query().Get("nodegroup")
	// // desiredSizeStr := r.URL.Query().Get("nodecount")
	// // http://192.168.0.51:4885/apis/eksinstancestart?region=ap-northeast-2&node=ip-172-31-58-160.ap-northeast-2.compute.internal
	// akid := "AKIAJGFO6OXHRN2H6DSA"                          //
	// secretkey := "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" //

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region), //
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}
	svc := ec2.New(sess)

	lists, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	var targetID string
	for _, v := range lists.Reservations {
		for _, e := range v.Instances {
			if *e.PrivateDnsName == nodename {
				targetID = *e.InstanceId
				break
			}
		}
	}

	asSvc := autoscaling.New(sess)

	var targetASG string
	instances, _ := asSvc.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{})
	for _, instance := range instances.AutoScalingInstances {
		if *instance.InstanceId == targetID {
			targetASG = *instance.AutoScalingGroupName
			break
		}
	}

	// fmt.Println("targetID", targetID, targetASG)

	input := &autoscaling.ExitStandbyInput{
		AutoScalingGroupName: aws.String(targetASG),
		InstanceIds: []*string{
			aws.String(targetID),
		},
	}

	_, err = asSvc.ExitStandby(input)
	if err != nil {
		errmsg := jsonErr{503, "failed", err.Error()}
		json.NewEncoder(w).Encode(errmsg)
	} else {
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(targetID),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StartInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StartInstances(input)
			if err != nil {
				errmsg := jsonErr{503, "failed", err.Error()}
				json.NewEncoder(w).Encode(errmsg)
			} else {
				fmt.Println("Success", result.StartingInstances)
				errmsg := jsonErr{200, "success", "vm instance starting"}
				json.NewEncoder(w).Encode(errmsg)
			}
		} else {
			errmsg := jsonErr{503, "failed", err.Error()}
			json.NewEncoder(w).Encode(errmsg)
		}
	}
}

//eks resource change
func ChangeEKSInstanceType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// region := r.URL.Query().Get("region")
	// nodename := r.URL.Query().Get("node")
	// instanceType := r.URL.Query().Get("type")

	// // http://192.168.0.51:4885/apis/changeekstype?region=ap-northeast-2&cluster=eks-cluster1&type=t3.nano&node=ip-172-31-58-160.ap-northeast-2.compute.internal
	// akid := "AKIAJGFO6OXHRN2H6DSA"                          //
	// secretkey := "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" //

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	region := data["region"].(string)
	instanceType := data["type"].(string)
	nodename := data["node"].(string)

	akid := data["akid"].(string)
	secretkey := data["secretKey"].(string)

	// fmt.Println("region : " + region)
	// fmt.Println("instanceType : " + instanceType)
	// fmt.Println("nodename : " + nodename)
	// fmt.Println("akid : " + akid)
	// fmt.Println("secretkey : " + secretkey)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region), //
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}

	svc := ec2.New(sess)
	lists, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	var targetID string
	for _, v := range lists.Reservations {
		for _, e := range v.Instances {
			if *e.PrivateDnsName == nodename {
				targetID = *e.InstanceId
				break
			}
		}
	}

	asSvc := autoscaling.New(sess)
	var targetASG string
	instances, _ := asSvc.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{})
	for _, instance := range instances.AutoScalingInstances {
		if *instance.InstanceId == targetID {
			targetASG = *instance.AutoScalingGroupName
			break
		}
	}

	// enter standby
	// fmt.Println("targetID", targetID, targetASG)
	input := &autoscaling.EnterStandbyInput{
		AutoScalingGroupName: aws.String(targetASG),
		InstanceIds: []*string{
			aws.String(targetID),
		},
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	}

	result, err := asSvc.EnterStandby(input)
	fmt.Println(result)
	if err != nil {
		errmsg := jsonErr{503, "failed", err.Error()}
		json.NewEncoder(w).Encode(errmsg)
		return
	} else {
		// stop instance
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(targetID),
			},
			DryRun: aws.Bool(true),
		}
		_, err = svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			_, err = svc.StopInstances(input)
			if err != nil {
				errmsg := jsonErr{503, "failed", err.Error()}
				json.NewEncoder(w).Encode(errmsg)
			} else {
				// fmt.Println("Success", result.StoppingInstances)
				fmt.Println("stop instance success")
			}
		} else {
			errmsg := jsonErr{503, "failed", err.Error()}
			json.NewEncoder(w).Encode(errmsg)
		}
	}

	for i := 0; i < 100; i++ {
		lists, _ = svc.DescribeInstances(&ec2.DescribeInstancesInput{})
		ck := "n"
		for _, v := range lists.Reservations {
			for _, e := range v.Instances {
				if *e.InstanceId == targetID {
					fmt.Println(*e.State.Name)
					if *e.State.Name == "stopped" {
						ck = "y"
						break
					}
				}
			}
		}
		if ck == "y" {
			break
		}
		time.Sleep(time.Second * 3)
	}

	// change instance type
	res, err := svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(targetID),
		InstanceType: &ec2.AttributeValue{
			Value: aws.String(instanceType),
		},
	})
	fmt.Println(res)

	if err != nil {
		fmt.Println(err)
		errmsg := jsonErr{503, "failed", err.Error()}
		json.NewEncoder(w).Encode(errmsg)
	} else {
		fmt.Println("change instance type success")
		// exit standby
		input := &autoscaling.ExitStandbyInput{
			AutoScalingGroupName: aws.String(targetASG),
			InstanceIds: []*string{
				aws.String(targetID),
			},
		}
		_, err = asSvc.ExitStandby(input)
		if err != nil {
			errmsg := jsonErr{503, "failed", err.Error()}
			json.NewEncoder(w).Encode(errmsg)
		} else {
			// start instance
			input := &ec2.StartInstancesInput{
				InstanceIds: []*string{
					aws.String(targetID),
				},
				DryRun: aws.Bool(true),
			}
			_, err = svc.StartInstances(input)
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "DryRunOperation" {
				input.DryRun = aws.Bool(false)
				_, err = svc.StartInstances(input)
				if err != nil {
					errmsg := jsonErr{503, "failed", err.Error()}
					json.NewEncoder(w).Encode(errmsg)
				} else {
					// fmt.Println("Success", result.StartingInstances)
					fmt.Println("start instacne success")
					successmsg := jsonErr{200, "success", "vm instance type updated"}
					json.NewEncoder(w).Encode(successmsg)
				}
			} else {
				errmsg := jsonErr{503, "failed", err.Error()}
				json.NewEncoder(w).Encode(errmsg)
			}
		}
	}

	// json.NewEncoder(w).Encode(dng.Nodegroup.Resources.AutoScalingGroups)
}
