package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	clusterPtr     = flag.String("cluster", "dev", "cluster to get list services and versions")
	profilePtr     = flag.String("profile", "", "set AWS profile, default will discover default profile")
	diffClusterPtr = flag.String("diffCluster", "", "cluster to compare difference")
	diffProfilePtr = flag.String("diffProfile", "", "override AWS profile for '-diffCluster' only if a cluster exists in a different AWS account")
	versionPtr     = flag.Bool("version", false, "show current version")
)

func main() {
	flag.Parse()

	if *versionPtr {
		fmt.Println("v1.2.0")
		os.Exit(0)
	}

	if *diffProfilePtr == "" {
		*diffProfilePtr = *profilePtr
	}

	serviceMap := getServiceVersions(*clusterPtr, *profilePtr)

	if *diffClusterPtr != "" {
		diffServiceMap := getServiceVersions(*diffClusterPtr, *diffProfilePtr)

		fmt.Println(strings.ToUpper(*clusterPtr), " => ", strings.ToUpper(*diffClusterPtr))
		printDiff(serviceMap, diffServiceMap)
	} else {
		fmt.Println(strings.ToUpper(*clusterPtr))
		printMap(serviceMap)
	}
}

func printDiff(x map[string]string, y map[string]string) {
	emptyStr := ""

	allKeys := removeDuplicates(append(getMapKeys(x), getMapKeys(y)...))
	sort.Strings(allKeys)

	for _, key := range allKeys {

		if _, ok := x[key]; !ok {
			x[key] = emptyStr
		}

		if _, ok := y[key]; !ok {
			y[key] = emptyStr
		}

		if x[key] == emptyStr && y[key] != emptyStr {
			color.Red(fmt.Sprint("-    ", key, ": ", y[key]))
		} else if x[key] != emptyStr && y[key] == emptyStr {
			color.Green(fmt.Sprint("+    ", key, ": ", x[key]))
		} else if x[key] != y[key] {
			color.Yellow(fmt.Sprint("~    ", key, ": ", x[key], " => ", y[key]))
		} else {
			color.White(fmt.Sprint("     ", key, ": ", x[key], " = ", y[key]))
			//fmt.Printf("     ", key, ": ", x[key], " = ", y[key])
		}
	}
}

func printMap(m map[string]string) {
	for key, value := range m {
		fmt.Println("\t"+key+":", value)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getMapKeys(x map[string]string) []string {
	i := 0
	keys := make([]string, len(x))
	for key := range x {
		keys[i] = key
		i++
	}
	return keys
}

func getServiceVersions(cluster string, profile string) map[string]string {
	resultMap := make(map[string]string)

	sessOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if profile != "" {
		fmt.Println("Using Profile", profile, "for cluster", cluster)
		sessOpts.Profile = profile
	}

	sess := session.Must(session.NewSessionWithOptions(sessOpts))

	svcECS := ecs.New(sess)

	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	result, err := svcECS.ListServices(input)
	check(err)

	for _, service := range result.ServiceArns {
		serviceOpts := ecs.DescribeServicesInput{
			Cluster:  aws.String(fmt.Sprintf(cluster)),
			Services: aws.StringSlice([]string{*service}),
		}
		serviceResponse, err := svcECS.DescribeServices(&serviceOpts)
		check(err)

		serviceName := *serviceResponse.Services[0].ServiceName

		taskDefOpts := ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(*serviceResponse.Services[0].TaskDefinition),
		}

		taskResponse, err := svcECS.DescribeTaskDefinition(&taskDefOpts)

		dockerImage := *taskResponse.TaskDefinition.ContainerDefinitions[0].Image

		resultMap[serviceName] = strings.Split(dockerImage, ":")[1]

	}

	return resultMap
}

func removeDuplicates(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
