package cloud

import (
	"log"

	"portal-api-server/db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
	"time"
)

func AddNode(nodenm string, aKey string, sKey string) AddNodeResult {

	var r = AddNodeResult{}

	akid := aKey
	secretkey := sKey

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	// Create EC2 service client
	svc := ec2.New(sess)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-0e077dbfdc14f6e35"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		r = AddNodeResult{"Could not create instance", "", ""}
		return r
	}

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(nodenm),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		r = AddNodeResult{"Could not create instance", *runResult.Instances[0].InstanceId, errtag.Error()}
		return r
	}

	fmt.Println("Successfully tagged instance")
	r = AddNodeResult{"Created instance", *runResult.Instances[0].InstanceId, ""}
	return r
}

func GetNodeState(instanceId *string, nodenm string, cluster string, aKey string, sKey string) {

	akid := aKey
	secretkey := sKey

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	// Create EC2 service client
	ec2Svc := ec2.New(sess)
	iid := []*string{instanceId}

	if err != nil {
		fmt.Println("Error", err)
	} else {
		for i := 0; i <= 120; i++ {
			result, errr := ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: iid,
			})
			status := *result.Reservations[0].Instances[0].State.Name
			publicIPAddress := ""
			provider := "aws"
			if errr != nil {
				fmt.Println("GetNodeState_Error", errr)
			} else {
				if status == "running" {
					publicIPAddress = *result.Reservations[0].Instances[0].PublicIpAddress
					db.InsertReadyNode(cluster, nodenm, publicIPAddress, status, provider)
					break
				} else {
					publicIPAddress = ""
					db.InsertReadyNode(cluster, nodenm, publicIPAddress, status, provider)
				}
			}
			time.Sleep(time.Second * 5)
		}

	}

}
