package lambda

import (
	"context"
	"fmt"

	real "github.com/aws/aws-sdk-go/service/lambda"
)

// FuncHandler is a longest function as the input for lambda.
// present for func (context.Context, TIn) (TOut, error)
// If you do not need to use input or output, just ignore it
type FuncHandler func(ctx context.Context, input []byte) ([]byte, error)

type Fake struct {
	mpLambdaFunc map[string]FuncHandler
}

type Option func(f *Fake)

func NewFakeLambda(opts ...Option) *Fake {
	f := &Fake{
		mpLambdaFunc: make(map[string]FuncHandler),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithFunction(name string, handler FuncHandler) Option {
	return func(f *Fake) {
		f.mpLambdaFunc[name] = handler
	}
}

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
