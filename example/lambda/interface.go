package main

import (
	"context"
	"sync"

	fake "github.com/NextOption/aws/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	real "github.com/aws/aws-sdk-go/service/lambda"
)

type Serverless interface {
	Invoke(in *real.InvokeInput) (out *real.InvokeOutput, err error)
	InvokeWithContext(ctx context.Context, in *real.InvokeInput, opts ...request.Option) (*real.InvokeOutput, error)
}

var (
	underlineLambda Serverless
	one             sync.Once
)

func GetServerless() Serverless {
	return underlineLambda
}

func InitFakeLambda(opts ...fake.Option) {
	one.Do(func() {
		underlineLambda = fake.NewFakeLambda(opts...)
	})
}

func InitAWSLambda() {
	one.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		client := real.New(sess, &aws.Config{Region: aws.String("us-west-2")})
		underlineLambda = client
	})
}
