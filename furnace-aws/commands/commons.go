package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/fatih/color"
	"github.com/go-furnace/go-furnace/handle"
)

func gatherParameters(source *os.File, params *cloudformation.ValidateTemplateOutput) []cloudformation.Parameter {
	var stackParameters []cloudformation.Parameter
	defaultValue := color.New(color.FgHiBlack, color.Italic).SprintFunc()
	log.Println("Gathering parameters.")
	if params == nil {
		return stackParameters
	}
	for _, v := range params.Parameters {
		var param cloudformation.Parameter
		fmt.Printf("%s - '%s'(%s):", aws.StringValue(v.Description), keyName(aws.StringValue(v.ParameterKey)), defaultValue(aws.StringValue(v.DefaultValue)))
		text := readInputFrom(source)
		param.ParameterKey = v.ParameterKey
		text = strings.Trim(text, "\n")
		if len(text) > 0 {
			param.ParameterValue = aws.String(text)
		} else {
			param.ParameterValue = v.DefaultValue
		}
		stackParameters = append(stackParameters, param)
	}
	return stackParameters
}

func readInputFrom(source *os.File) string {
	reader := bufio.NewReader(source)
	text, _ := reader.ReadString('\n')
	return text
}

func (cf *CFClient) describeStacks(descStackInput *cloudformation.DescribeStacksInput) *cloudformation.DescribeStacksOutput {
	req := cf.Client.DescribeStacksRequest(descStackInput)
	descResp, err := req.Send()
	handle.Error(err)
	return descResp
}

func (cf *CFClient) validateTemplate(template []byte) *cloudformation.ValidateTemplateOutput {
	log.Println("Validating template.")
	validateParams := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(string(template)),
	}
	req := cf.Client.ValidateTemplateRequest(validateParams)
	resp, err := req.Send()
	handle.Error(err)
	return resp
}
