package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stefan79/gadgeto/pkg/bootstrap"
	"github.com/stefan79/gadgeto/pkg/resources/aws/apigw"
	"github.com/stefan79/gadgeto/pkg/resources/aws/gs3"
	"github.com/stefan79/gadgeto/pkg/triggers/route"
)

type (
	Setup struct {
		MyS3Bucket gs3.Client
	}
	Deploy struct {
		MyAPIGW apigw.APIGatewayClient
	}
)

type Input struct {
	Key  string
	Body string
}

type Output struct {
	Message string
}

func (s *Setup) CreateCustomer(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	name := request.QueryStringParameters["name"]
	message := fmt.Sprintf("Hello %s", name)
	err := s.MyS3Bucket.WriteToObject(name, []byte(message))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}

func main() {
	ctx := bootstrap.NewContext()

	deployment := &Deploy{
		MyAPIGW: apigw.
			NewApiGatewayClient("demo", ctx).
			AddTag("stefan", "gadget").
			Build(),
	}

	setup := &Setup{
		MyS3Bucket: gs3.
			S3(ctx, "myBucket").
			WithBucketName("stefansiprell1979test").
			Build(),
	}

	route.NewTrigger("putCustomer", deployment.MyAPIGW).
		WithKey(route.POST, "/customers").
		Build().
		Handle(setup.CreateCustomer)

	ctx.Complete()
}
