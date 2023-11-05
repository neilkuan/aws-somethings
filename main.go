package main

import (
	"aws-somethings/iam"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("have something wrong...:  %v", err)
	}
}

func GetIamConfigure() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	checkError(err)
	return cfg
}

func main() {
	defConfg := GetIamConfigure()
	newPolicyList := iam.ListIAMPolicy(defConfg)
	oldPolicy, err := os.ReadFile("policies.json")
	checkError(err)

	var oldPolicyList []*string
	err = json.Unmarshal(oldPolicy, &oldPolicyList)
	checkError(err)

	enCodeList, err := json.MarshalIndent(newPolicyList, "", "    ")
	checkError(err)

	updateFile, err := os.Create("update.txt")
	checkError(err)

	if !(string(oldPolicy) == string(enCodeList)) {
		log.Println("Start to Update policies.json...")
		err = os.WriteFile("policies.json", enCodeList, 0644)
		checkError(err)
		log.Println("Update policies.json successfully!")

		log.Println("Updated update.txt...")
		updateFile.WriteString("Y")

	} else {
		log.Println("Not need to update update.txt...")
		updateFile.WriteString("N")
	}
}
