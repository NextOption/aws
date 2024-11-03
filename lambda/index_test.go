package lambda

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/lambda"
	"reflect"
	"testing"
)

func Test_isAsyncInvoke(t *testing.T) {
	type args struct {
		input lambda.InvokeInput
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "InvocationType is Event",
			args: args{
				input: lambda.InvokeInput{
					InvocationType: aws.String(lambda.InvocationTypeEvent),
				},
			},
			want: true,
		},
		{
			name: "InvocationType is RequestResponse",
			args: args{
				input: lambda.InvokeInput{
					InvocationType: aws.String(lambda.InvocationTypeRequestResponse),
				},
			},
			want: false,
		},
		{
			name: "InvocationType is nil",
			args: args{
				input: lambda.InvokeInput{
					InvocationType: nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAsyncInvoke(tt.args.input); got != tt.want {
				t.Errorf("isAsyncInvoke() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFakeLambda(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want *Fake
	}{
		{
			name: "No options",
			opts: nil,
			want: &Fake{
				mpLambdaFunc: make(map[string]FuncHandler),
			},
		},
		{
			name: "With one function",
			opts: []Option{
				WithFunction("testFunc", func(ctx context.Context, input []byte) ([]byte, error) {
					return []byte("output"), nil
				}),
			},
			want: &Fake{
				mpLambdaFunc: map[string]FuncHandler{
					"testFunc": func(ctx context.Context, input []byte) ([]byte, error) {
						return []byte("output"), nil
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFakeLambda(tt.opts...)
			if len(got.mpLambdaFunc) != len(tt.want.mpLambdaFunc) {
				t.Errorf("NewFakeLambda() mpLambdaFunc length = %v, want %v", len(got.mpLambdaFunc), len(tt.want.mpLambdaFunc))
				return
			}
			for key, wantFunc := range tt.want.mpLambdaFunc {
				gotFunc, ok := got.mpLambdaFunc[key]
				if !ok {
					t.Errorf("NewFakeLambda() missing function %v", key)
					continue
				}
				gotOutput, gotErr := gotFunc(context.TODO(), []byte("input"))
				wantOutput, wantErr := wantFunc(context.TODO(), []byte("input"))
				if !reflect.DeepEqual(gotOutput, wantOutput) || !reflect.DeepEqual(gotErr, wantErr) {
					t.Errorf("NewFakeLambda() function %v = (%v, %v), want (%v, %v)", key, gotOutput, gotErr, wantOutput, wantErr)
				}
			}
		})
	}
}

func TestFake_validateInput(t *testing.T) {
	type fields struct {
		mpLambdaFunc map[string]FuncHandler
	}
	type args struct {
		in *lambda.InvokeInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "No input",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{},
			},
			args: args{
				in: nil,
			},
			wantErr: true,
		},
		{
			name: "No function name",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{},
			},
			args: args{
				in: &lambda.InvokeInput{},
			},
			wantErr: true,
		},
		{
			name: "Function not found",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{},
			},
			args: args{
				in: &lambda.InvokeInput{
					FunctionName: aws.String("nonExistentFunc"),
				},
			},
			wantErr: true,
		},
		{
			name: "Valid input",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{
					"testFunc": func(ctx context.Context, input []byte) ([]byte, error) {
						return []byte("output"), nil
					},
				},
			},
			args: args{
				in: &lambda.InvokeInput{
					FunctionName: aws.String("testFunc"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fake{
				mpLambdaFunc: tt.fields.mpLambdaFunc,
			}
			if err := f.validateInput(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("validateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFake_InvokeWithContext(t *testing.T) {
	type fields struct {
		mpLambdaFunc map[string]FuncHandler
	}
	type args struct {
		ctx  context.Context
		in   *lambda.InvokeInput
		opts []request.Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *lambda.InvokeOutput
		wantErr bool
	}{
		{
			name: "No input",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{},
			},
			args: args{
				ctx: context.TODO(),
				in:  nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Function not found",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{},
			},
			args: args{
				ctx: context.TODO(),
				in: &lambda.InvokeInput{
					FunctionName: aws.String("nonExistentFunc"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Valid input, sync invoke",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{
					"testFunc": func(ctx context.Context, input []byte) ([]byte, error) {
						return []byte("output"), nil
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				in: &lambda.InvokeInput{
					FunctionName: aws.String("testFunc"),
					Payload:      []byte("input"),
				},
			},
			want: &lambda.InvokeOutput{
				Payload: []byte("output"),
			},
			wantErr: false,
		},
		{
			name: "Valid input, async invoke",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{
					"testFunc": func(ctx context.Context, input []byte) ([]byte, error) {
						return []byte("output"), nil
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				in: &lambda.InvokeInput{
					FunctionName:   aws.String("testFunc"),
					Payload:        []byte("input"),
					InvocationType: aws.String(lambda.InvocationTypeEvent),
				},
			},
			want: &lambda.InvokeOutput{
				StatusCode: aws.Int64(202),
			},
			wantErr: false,
		},
		{
			name: "Handler returns error",
			fields: fields{
				mpLambdaFunc: map[string]FuncHandler{
					"errorFunc": func(ctx context.Context, input []byte) ([]byte, error) {
						return nil, fmt.Errorf("handler error")
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				in: &lambda.InvokeInput{
					FunctionName: aws.String("errorFunc"),
					Payload:      []byte("input"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fake{
				mpLambdaFunc: tt.fields.mpLambdaFunc,
			}
			got, err := f.InvokeWithContext(tt.args.ctx, tt.args.in, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvokeWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvokeWithContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}
