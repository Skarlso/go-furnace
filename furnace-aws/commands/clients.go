package commands

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/codedeployiface"
	"github.com/aws/aws-sdk-go-v2/service/iam/iamiface"
	"github.com/go-furnace/go-furnace/handle"
	"log"
)

// CFClient abstraction for cloudFormation client.
type CFClient struct {
	Client cloudformationiface.ClientAPI
}

func (cf *CFClient) describeStacks(descStackInput *cloudformation.DescribeStacksInput) *cloudformation.DescribeStacksResponse {
	req := cf.Client.DescribeStacksRequest(descStackInput)
	descResp, err := req.Send(context.Background())
	handle.Error(err)
	return descResp
}

func (cf *CFClient) validateTemplate(template []byte) *cloudformation.ValidateTemplateResponse {
	log.Println("Validating template.")
	validateParams := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(string(template)),
	}
	req := cf.Client.ValidateTemplateRequest(validateParams)
	resp, err := req.Send(context.Background())
	handle.Error(err)
	return resp
}

// CDClient abstraction for cloudFormation client.
type CDClient struct {
	Client codedeployiface.ClientAPI
}

// IAMClient abstraction for cloudFormation client.
type IAMClient struct {
	Client iamiface.ClientAPI
}

// NotEmptyStack test structs which defines a non-empty stack.
var NotEmptyStack = cloudformation.DescribeStacksOutput{
	Stacks: []cloudformation.Stack{
		{
			StackName:   aws.String("TestStack"),
			StackStatus: cloudformation.StackStatusCreateComplete,
		},
	},
}
