package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

type MyEvent struct {
	Hello string `json:"hello"`
}

func handler(ctx context.Context, name MyEvent) (string, error) {
	log.Printf("name: %v", name)
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if len(ip) == 0 {
		return "", ErrNoIP
	}

	return fmt.Sprintf("Hello, %v", string(ip)), nil
}

func main() {
	lambda.Start(handler)
}
