package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	key                    = flag.String("key", "", "")
	secret                 = flag.String("secret", "", "")
	region                 = flag.String("region", "", "")
	clusterName            = flag.String("cluster-name", "", "")
	serviceName            = flag.String("service-name", "", "")
	taskDefinitionName     = flag.String("task-definition-name", "", "")
	taskDefinitionFileName = flag.String("task-definition-file", "", "")
	desiredTasksCount      = flag.Int64("desired-tasks-count", 0, "")
)

func main() {
	flag.Parse()
	fmt.Println("Hello world!")

	log.Println("Step 1: Configuring AWS")
	ecsCli := ecs.New(session.New(aws.NewConfig().
		WithRegion("eu-west-1").
		WithCredentials(credentials.NewStaticCredentials(*key, *secret, "foo"))))
	log.Println("Configuring AWS succeeded")

	log.Println("Step 2: Check ECS cluster")
	describeClusterOut, err := ecsCli.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{clusterName},
	})
	if err != nil {
		log.Fatal(err)
	}
	if ff := describeClusterOut.Failures; len(ff) > 0 {
		log.Fatalf("Cluster resource '%s' failed: %v", *ff[0].Arn, *ff[0].Reason)
	}
	log.Printf("Checking cluster '%s' succeeded", *clusterName)

	log.Println("Step 3: Check ECS Service")
	describeServicesOut, err := ecsCli.DescribeServices(&ecs.DescribeServicesInput{
		Services: []*string{serviceName},
	})
	if err != nil {
		log.Fatal(err)
	}
	origRunningCount := describeServicesOut.Services[0].RunningCount
	log.Printf("Checking service '%s' succeeded (%d tasks running)", *serviceName, *origRunningCount)

	log.Println("Step 4: Register New Task Definition")
	var cd []*ecs.ContainerDefinition
	taskDefinitionFile, err := os.Open(*taskDefinitionFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer taskDefinitionFile.Close()

	err = json.NewDecoder(taskDefinitionFile).Decode(&cd)
	if err != nil {
		log.Fatal(err)
	}

	registerTaskDefinitionOut, err := ecsCli.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: cd,
		Family:               taskDefinitionName,
	})
	if err != nil {
		log.Fatal(err)
	}

	if *registerTaskDefinitionOut.TaskDefinition.Status == "INACTIVE" {
		log.Fatalf("Task definition (%s) is inactive", *registerTaskDefinitionOut.TaskDefinition.TaskDefinitionArn)
	}
	taskDefinitionArn := registerTaskDefinitionOut.TaskDefinition.TaskDefinitionArn
	log.Printf("Registering task definition '%s' succeeded", *taskDefinitionArn)

	_, err = ecsCli.UpdateService(&ecs.UpdateServiceInput{
		DesiredCount:   desiredTasksCount,
		Service:        serviceName,
		TaskDefinition: taskDefinitionArn,
	})
	if err != nil {
		log.Fatal(err)
	}
}
