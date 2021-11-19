package iam

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("have something wrong...:  %v", err)
	}
}

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
