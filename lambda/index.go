package lambda

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/request"
	real "github.com/aws/aws-sdk-go/service/lambda"
)

// isAsyncInvoke checks if the input is async invoke
func isAsyncInvoke(input real.InvokeInput) bool {
	return input.InvocationType != nil && *input.InvocationType == real.InvocationTypeEvent
}

// FuncHandler is a longest function as the input for lambda.
// present for func (context.Context, TIn) (TOut, error)
// If you do not need to use input or output, just ignore it
type FuncHandler func(ctx context.Context, input []byte) ([]byte, error)

// Fake presents a fake lambda that keep mapping between function name and handler
type Fake struct {
	mpLambdaFunc map[string]FuncHandler
}

// Option add option for Fake
type Option func(f *Fake)

// WithFunction register a function with Fake
func WithFunction(name string, handler FuncHandler) Option {
	return func(f *Fake) {
		f.mpLambdaFunc[name] = handler
	}
}

// NewFakeLambda creates a new fake lambda
func NewFakeLambda(opts ...Option) *Fake {
	f := &Fake{
		mpLambdaFunc: make(map[string]FuncHandler),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// validateInput validates the input for lambda
func (f *Fake) validateInput(in *real.InvokeInput) error {
	if in == nil {
		return ErrNoInput
	}
	if in.FunctionName == nil {
		return ErrNoFunctionName
	}
	if _, ok := f.mpLambdaFunc[*in.FunctionName]; !ok {
		return fmt.Errorf("not found function: %s", *in.FunctionName)
	}
	return nil
}

// Invoke invokes the lambda function-same signature as real lambda
func (f *Fake) Invoke(in *real.InvokeInput) (out *real.InvokeOutput, err error) {
	return f.InvokeWithContext(context.TODO(), in)
}

// InvokeWithContext invokes the lambda function-same signature as real lambda
func (f *Fake) InvokeWithContext(ctx context.Context, in *real.InvokeInput, opts ...request.Option) (*real.InvokeOutput, error) {
	if err := f.validateInput(in); err != nil {
		return nil, err
	}
	handler := f.mpLambdaFunc[*in.FunctionName]
	if isAsyncInvoke(*in) {
		go func() {
			_, _ = handler(context.TODO(), in.Payload)
		}()
		return &real.InvokeOutput{
			StatusCode: &StatusCode202,
		}, nil
	}
	payload, err := handler(ctx, in.Payload)
	if err != nil {
		return nil, err
	}
	return &real.InvokeOutput{
		Payload: payload,
	}, nil
}

// InvokeAsync, InvokeAsyncWithContext invokes the lambda function-same signature as real lambda
// Deprecated: InvokeAsync has been deprecated in aws api https://docs.aws.amazon.com/lambda/latest/api/API_InvokeAsync.html
// So we skip it. Refer to Invoke to invoke asynchronously
// func (f *Fake) InvokeAsync(input *real.InvokeAsyncInput) (*real.InvokeAsyncOutput, error){}
// }

// Deprecated https://docs.aws.amazon.com/lambda/latest/api/API_InvokeAsync.html
// func (f *Fake ) InvokeAsyncWithContext(aws.Context, *lambda.InvokeAsyncInput, ...request.Option) (*lambda.InvokeAsyncOutput, error){
// }
