package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func ListIAMPolicy(cfg aws.Config) []*string {
	iamsvc := iam.NewFromConfig(cfg)
	plist := []*string{}
	truncatedListing := true
	var resp *iam.ListPoliciesOutput
	var err error
	params := &iam.ListPoliciesInput{
		Scope:    types.PolicyScopeTypeAws,
		MaxItems: aws.Int32(1000),
	}
	for truncatedListing {
		resp, err = iamsvc.ListPolicies(context.TODO(), params)
		checkError(err)
		for _, policy := range resp.Policies {
			plist = append(plist, policy.PolicyName)
		}
		params.Marker = resp.Marker
		truncatedListing = resp.IsTruncated
	}

	return plist
}

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
	newPolicyList := ListIAMPolicy(defConfg)
	oldPolicy, err := ioutil.ReadFile("policies.json")
	checkError(err)

	var oldPolicyList []*string
	err = json.Unmarshal(oldPolicy, &oldPolicyList)
	checkError(err)

	enCodeList, err := json.MarshalIndent(newPolicyList, "", "    ")
	checkError(err)

	policyNewTemp, err := os.Create("policy-new-temp.json")
	checkError(err)
	_, err = policyNewTemp.Write(enCodeList)
	checkError(err)
	newPolicy, err := ioutil.ReadFile("policy-new-temp.json")
	checkError(err)

	updateFile, err := os.Create("update.txt")
	checkError(err)

	if !bytes.Equal(oldPolicy, newPolicy) {
		log.Println("Start to Update policies.json...")
		err = ioutil.WriteFile("policies.json", enCodeList, 0644)
		checkError(err)
		log.Println("Update policies.json successfully!")

		log.Println("Updated update.txt...")
		updateFile.WriteString("Y")

	} else {
		log.Println("Not need to update update.txt...")
		updateFile.WriteString("N")
	}

	log.Println("start to remove temp file...")
	os.Remove("policy-new-temp.json")
}
