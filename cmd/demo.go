package main

import (
	"context"

	"github.com/stefan79/gadgeto/pkg/bootstrap"
	"github.com/stefan79/gadgeto/pkg/resources/aws/gs3"
	"github.com/stefan79/gadgeto/pkg/triggers/native"
)

type (
	Setup struct {
		MyS3Bucket gs3.Client
	}
)

/*
func (s *Setup) handleCreateCustomerCall(request triggers.Request, response triggers.Response) error {
	key, ok := request.QueryParams["key"]
	var err error
	if !ok {
		response.ResponseCode = 400
	} else {
		err = s.MyS3Bucket.WriteToObject(key, request.Body)
	}
	return err
}
*/

func (s *Setup) handleNativeCall(ctx context.Context, input Input) (Output, error) {
	err := s.MyS3Bucket.WriteToObject(input.Key, []byte(input.Body))
	return Output{
		Message: "OK",
	}, err
}

type Input struct {
	Key  string
	Body string
}

type Output struct {
	Message string
}

func main() {
	ctx := bootstrap.NewContext()
	/*NewContext().
	WithAppName("demo").
	WithTags("stefan", "gadget").
	WithPermission("s3:PutObject", "*").
	Build()*/

	setup := &Setup{
		MyS3Bucket: gs3.
			S3(ctx, "myBucket").
			WithBucketName("stefansiprell1979test").
			Build(),
	}
	native.
		NewNativeTrigger[Input, Output]("mainTrigger", ctx).
		Build().
		Handle(setup.handleNativeCall)

		//apigw.ApiGateway("CreateCustomer").WithMethod("POST").Build(ctx).Handle(setup.handleCreateCustomerCall)

	ctx.Complete()
}
