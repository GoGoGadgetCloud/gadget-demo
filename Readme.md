# GoGo Gadget

This is a demo project showing off how we could mix and match deployment code with business code. Goal is to bring the new kind of tooling to the Go World, as Go is an excellent language predestined for serverless workloads.

## The Code

### Setting up Resources

GoGo Gadget is completely unbiased how you structure your code. The only requisite is that the main method will be invoking the creators of all of the resources.

```
apigw.
    NewApiGatewayClient("demo", ctx).
    AddTag("stefan", "gadget").
    Build()

gs3.
    S3(ctx, "myBucket").
	WithBucketName("stefansiprell1979test").
	Build()
```
The resources are configured with a fluent API and with full compiler support. The API should be targeted towards develoeprs. So we would not be denfining the authentication via multiple strings - like cloudformation requires - but by dropping in special interfaces, and let the code wire everything up.

This will limit the flexibility but we can make special implementations of resources (like a separate impl suited for websockets) similar the way you can make level 3 constructs in CDK code to optimimze for your local requirements.

Also we should add special "helper classes" as resources, e.g. a logger which outputs structured JSON messages on the server, but outputs colored console output on the local machine. 

### Settin Up Triggers
At the moment we support native triggers (arbitrary struct mapping for request and response events) and apigateways. Next up is to implement an S3 trigger to show how you can quickly chain together complex applications with simple code.

Triggers use the new generics extensively, and will make it impossible to write go code which does not use matching interfaces / structs.

```
	route.NewTrigger("putCustomer", deployment.MyAPIGW).
		WithKey(route.POST, "/customers").
		Build().
		Handle(setup.CreateCustomer)

```

Again GoGo Gadget is unbiased towards the actual implementation. You can mix and match your own code, middleware, whatever. In this case we use our resource to create a s3 entry.

```
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
```

### The boilerplate 

```
	ctx := bootstrap.NewContext()
	ctx.Complete()
```

This is it.  The code does not interweave with your application logic. We merely need to let gogo gadget when it can do something. 

## The magic

### Overview

GoGo Gadget works with three components:

* gadget library: the library is embedded into the client code and "hijacks" the control depending on various parameters and will let the client do various things. Thhe library workd completely without code generation. 
* client: the client is a regular go project.
* cli: a small CLI application will drive the client controlled by the gadget library.

### Client Modes

The library can set the client into various modes, at the moment we support two:

* build mode: the library will create a complete cloud formation template for its command (command is something that `go build` creates and is mapped to a lambda function).
* run mode: the library detects an AWS runtime via environment settings and starts into "run mode".

THe non run modes are configured via command line interfaces and are quite comfortable, below is a sample of an incompletely configured build mode:

```
gadget-demo on  main [?] via 🐹 v1.21.5
➜ ./.gadget/staging/demo_local deployment generate
NAME:
   demo_local deployment generate

USAGE:
   demo_local deployment generate [command options] [arguments...]

OPTIONS:
   --template value     where the template should be generated to
   --handler value      handler to use in the function definition
   --application value  application prefix to use in the created resources
   --command value      command prefix to use in the created resources
   --s3bucket value     s3 bucket name to use in the function source
   --s3key value        s3 key to use in the function source
   --help, -h           show help
ERRO Failed to start in desired mode err="Required flags \"template, handler, application, command, s3bucket, s3key\" not set"
```

Further modes can be envisioned:
* local mode: Starting the application locally and exposing HTTPS API / Lambda Interfaces locally by pulling the remote environment settings to the developer machine.
* proxy mode: Similar to SST by deploying a generic stub and using the IOT-GW to relay requests/responses to/from the local machine.
* test mode: Running predefined tests locally 
* least-privilege-scan mode: Probably not a native mode in the client, but allowing the CLI to analyze the AST to determine the least amount of required AWS policies to run the code.

Using the plug-in mechanism of Go introduced a while a go, we can define modes into seperate binary plugins and let the cli add them to the client as required. This way we can reduce size ,complexity and security exposure for the final client.

### The CLI

The command line has various tasks it can complete:

#### Bootstrap

`gadget bootstrap` will setup the AWS Bucket (currently not prpoperly secured) and other resources required for deploying the applications.

#### Workspace

`gadget workspace` will setup the local gadget.yaml descritor. It closely resembles go `go work` works. The current file:

```
name: GadgetDemo
commands:
- name: demo
  path: cmd/demo.go
tags:
  createdBy: stefan.siprell@gmail.com
```

You can use multi commands, i.e. deploy multiple lambdas from the same go project. Exactly how you would work with regular go commands. The tags will be added to all resources (tagging is still buggy) and the cli will pull in additional tags, like the `go.mod` module path and version for CI/CD transparency.

#### Deploy

A deployment will do the following:
- compile the application for local execution
- x-compile the application for remote execution
- zip the remote file
- checksum the remote file (working on skipping unnecessary uploads)
- run the client in build mode 
- upload the remote file
- create / update the cloudformation setup (working on lambda updates and skipping cloudformation)

It will do the command tasks seperately and merge the generated cloudformation file.

#### Staging

Is required but not even envisioned yet. Ideally we want to use the native mechanisms but this needs to be planned out.




