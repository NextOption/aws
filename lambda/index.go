package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/request"

	real "github.com/aws/aws-sdk-go/service/lambda"
)

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
	if err := f.validateInput(in); err != nil {
		return nil, err
	}
	handler := f.mpLambdaFunc[*in.FunctionName]
	payload, err := handler(context.TODO(), in.Payload)
	if err != nil {
		return nil, err
	}
	return &real.InvokeOutput{
		Payload: payload,
	}, nil
}

// InvokeWithContext invokes the lambda function-same signature as real lambda
func (f *Fake) InvokeWithContext(ctx context.Context, in *real.InvokeInput, opts ...request.Option) (*real.InvokeOutput, error) {
	if err := f.validateInput(in); err != nil {
		return nil, err
	}
	handler := f.mpLambdaFunc[*in.FunctionName]
	payload, err := handler(ctx, in.Payload)
	if err != nil {
		return nil, err
	}
	return &real.InvokeOutput{
		Payload: payload,
	}, nil
}
