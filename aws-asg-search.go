package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func getWorkerNodesASGName(sess *session.Session, stackID string) (string, int64, error) {
	fmt.Println("getWorkerNodesASGName: Locating worker node auto-scaling group for", "stackID = ", stackID)
	svc := cloudformation.New(sess)

	input := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(stackID),
	}
	result, err := svc.ListStackResources(input)

	// Search for the k8s-node autoscaling group in stack resources
	for _, asc := range result.StackResourceSummaries {
		if strings.Contains(*asc.ResourceType, "AWS::AutoScaling::AutoScalingGroup") && strings.Contains(*asc.PhysicalResourceId, "K8sNodeGroup") && strings.Contains(*asc.PhysicalResourceId, stackID) {
			fmt.Println("getWorkerNodesASGName: Found worker node auto-scaling group for", "stackID =", stackID, ", AutoScalingGroup =", *asc.PhysicalResourceId)

			// extract info of autoscaling group
			svcAsg := autoscaling.New(sess)
			inputAsg := &autoscaling.DescribeAutoScalingGroupsInput{
				AutoScalingGroupNames: []*string{
					aws.String(*asc.PhysicalResourceId),
				},
			}
			resultAsg, _ := svcAsg.DescribeAutoScalingGroups(inputAsg)
			fmt.Println("getWorkerNodesASGName: AutoScaling group details :", "stackID = ", stackID, ", AutoScalingGroupName =", *asc.PhysicalResourceId, ", DesiredCapacity =", *resultAsg.AutoScalingGroups[0].DesiredCapacity)
			return *asc.PhysicalResourceId, *resultAsg.AutoScalingGroups[0].DesiredCapacity, nil
		}
	}

	fmt.Println("getWorkerNodesASGName: Failed to find worker node auto-scaling group for", "stackID", stackID)
	return "", 0, err
}

func main() {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	getWorkerNodesASGName(sess, "stack-6kn0y82")

}
