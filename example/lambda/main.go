package main

import (
	"context"
	"encoding/json"
	"fmt"

	fake "github.com/NextOption/aws/lambda"
	real "github.com/aws/aws-sdk-go/service/lambda"
)

type Order struct {
	ID string `json:"id"`
}

type OrderProcessResult struct {
	OrderID string `json:"order_id"`
	Success bool   `json:"success"`
}

// handleFn1 is handler for lambda_fn_1
func handleFn1(ctx context.Context, input []byte) ([]byte, error) {
	// because you do not need input and output, so we just ignore it
	coreFn1()
	return []byte{}, nil
}

func coreFn1() {
	fmt.Println("I'm inside of coreFn1")
}

// handleFn2 is handler for lambda_fn_2
func handleFn2(ctx context.Context, input []byte) ([]byte, error) {
	order := Order{}
	err := json.Unmarshal(input, &order)
	if err != nil {
		fmt.Println("error when unmarshal input: ", err)
		return []byte{}, err
	}
	success, err := coreFn2(&order)
	result := OrderProcessResult{
		OrderID: order.ID,
		Success: success,
	}
	out, err := json.Marshal(&result)
	if err != nil {
		fmt.Println("error when marshal output: ", err)
		return []byte{}, err
	}
	return out, nil
}

func coreFn2(order *Order) (success bool, err error) {
	fmt.Println("I'm inside of coreFn2, receive: ", order)
	return true, nil
}

func main() {
	var (
		lambdaFn1Name = "lambda_fn_1"
		fn2Name       = "lambda_fn_2"
	)
	enableFakeLambda := true
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

	order := Order{
		ID: "uuid",
	}
	inputPayload, err := json.Marshal(&order)
	if err != nil {
		fmt.Println("error when marshal input: ", err)
		return
	}
	invokeOutput, err := GetServerless().Invoke(&real.InvokeInput{
		FunctionName: &fn2Name,
		Payload:      inputPayload,
	})
	if err != nil {
		fmt.Printf("error when invoke lambda 2: %v", err)
		return
	}
	if invokeOutput != nil && invokeOutput.Payload != nil {
		orderProcessResult := OrderProcessResult{}
		err = json.Unmarshal(invokeOutput.Payload, &orderProcessResult)
		if err != nil {
			fmt.Println("error when unmarshal output: ", err)
			return
		}
		fmt.Println("result: ", orderProcessResult)
	}

}

// real main for lambda will be like following
//func main() {
//	lambda.Start(handleFn1)
//}

//func main() {
//	lambda.Start(handleFn2)
//}
