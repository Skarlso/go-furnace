package commands

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	commander "github.com/Yitsushi/go-commander"
	"github.com/go-furnace/go-furnace/config"
	awsconfig "github.com/go-furnace/go-furnace/furnace-aws/config"
	"github.com/go-furnace/go-furnace/handle"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/cloudformationiface"
)

type fakeDeleteCFClient struct {
	cloudformationiface.ClientAPI
	stackname string
	err       error
}

func (fc *fakeDeleteCFClient) DeleteStackRequest(*cloudformation.DeleteStackInput) cloudformation.DeleteStackRequest {
	return cloudformation.DeleteStackRequest{
		Request: &aws.Request{
			Data:        &cloudformation.DeleteStackOutput{},
			HTTPRequest: new(http.Request),
		},
	}
}

func (fc *fakeDeleteCFClient) WaitUntilStackDeleteComplete(ctx context.Context, input *cloudformation.DescribeStacksInput, opts ...aws.WaiterOption) error {
	return fc.err
}

func TestDeleteProcedure(t *testing.T) {
	config.WAITFREQUENCY = 0
	client := new(CFClient)
	stackname := "ToDeleteStack"
	client.Client = &fakeDeleteCFClient{err: nil, stackname: stackname}
	deleteStack(stackname, client)
}

func TestDeleteExecute(t *testing.T) {
	config.WAITFREQUENCY = 0
	client := new(CFClient)
	stackname := "ToDeleteStack"
	client.Client = &fakeDeleteCFClient{err: nil, stackname: stackname}
	opts := &commander.CommandHelper{}
	d := Delete{
		client: client,
	}
	d.Execute(opts)
}

func TestDeleteExecuteWithExtraStack(t *testing.T) {
	config.WAITFREQUENCY = 0
	client := new(CFClient)
	stackname := "ToDeleteStack"
	client.Client = &fakeDeleteCFClient{err: nil, stackname: stackname}
	opts := &commander.CommandHelper{}
	opts.Args = append(opts.Args, "teststack")
	d := Delete{
		client: client,
	}
	d.Execute(opts)
	if awsconfig.Config.Main.Stackname != "MyStack" {
		t.Fatal("test did not load the file requested.")
	}
}

func TestDeleteExecuteWithExtraStackNotFound(t *testing.T) {
	failed := false
	handle.LogFatalf = func(s string, a ...interface{}) {
		failed = true
	}
	config.WAITFREQUENCY = 0
	client := new(CFClient)
	stackname := "ToDeleteStack"
	client.Client = &fakeDeleteCFClient{err: nil, stackname: stackname}
	opts := &commander.CommandHelper{}
	opts.Args = append(opts.Args, "notfound")
	d := Delete{
		client: client,
	}
	d.Execute(opts)
	if !failed {
		t.Error("Expected outcome to fail. Did not fail.")
	}
}

func TestDeleteCreate(t *testing.T) {
	wrapper := NewDelete("furnace")
	if wrapper.Help.Arguments != "custom-config" ||
		!reflect.DeepEqual(wrapper.Help.Examples, []string{"", "custom-config"}) ||
		wrapper.Help.LongDescription != `Delete a stack with a given name.` ||
		wrapper.Help.ShortDescription != "Delete a stack" {
		t.Log(wrapper.Help.LongDescription)
		t.Log(wrapper.Help.ShortDescription)
		t.Log(wrapper.Help.Examples)
		t.Fatal("wrapper did not match with given params")
	}
}
