<p align="center">
<img src="https://raw.githubusercontent.com/NextOption/id/refs/heads/main/logo/NextOptionLogo.jpeg" alt="Next Option Logo">
<br/><br/>

## Next Option of AWS, LocalStack

_AWS did not support local development.  
LocalStack was born to solve this problem.  
But now they come with Premium Plan (-__-!) and provide cloud services.  
We are here with you, the next option of AWS, LocalStack._

## Feature

* Support for local development for RealAWS services.
* Become a Localstack alternative.
* Ready to live on production as a backup for RealAWS.

### Lambda

Just inject our fake lambda inside your code and you are ready to go.
It has the same interface as AWS Lambda but run in your code.
Only provide basic functions to Invoke now. It's enough for most of the cases that we just need to invoke lambda.

Supported Interface:

1. [x] `Invoke(input *InvokeInput) (*InvokeOutput, error)`
2. [x] `InvokeWithContext(aws.Context, *lambda.InvokeInput, ...request.Option) (*lambda.InvokeOutput, error)`
3. [x] `InvokeAsync` & `InvokeAsyncWithContext` deprecated by aws but supported by `Invoke` or `InvokeWithContext`
4. [ ] Other APIs are not supported yet.

Sample code:
```go
	if enableFakeLambda {
		InitFakeLambda(
			fake.WithFunction(lambdaFn1Name, handleFn1),
			fake.WithFunction(fn2Name, handleFn2),
		)
	} else {
		InitAWSLambda()
	}
	// invoke lambda_fn_1
	// lambda_fn_1 return nothing so we just skip it
	_, err := GetServerless().Invoke(&real.InvokeInput{
		FunctionName: &lambdaFn1Name,
	})
	if err != nil {
		fmt.Printf("error when invoke lambda 1: %v", err)
		return
	}
```

Check the [example/lambda](example/lambda) folder for more information.

### SQS (in progress)
Just try to implement simple queue by go channel.
