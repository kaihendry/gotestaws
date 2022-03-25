package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MyEvent struct {
	Hello string `json:"hello"`
}

// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

func handler(ctx context.Context, name MyEvent) (string, error) {
	log.Printf("name: %v", name)
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Non 200 Response found")
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	response, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String("ip-address"),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   bytes.NewReader(ip),
	})

	return fmt.Sprintf("Hello, %v", response), nil
}

func main() {
	lambda.Start(handler)
}
