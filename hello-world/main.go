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

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func StoreIpAddress(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

func handler(ctx context.Context, name MyEvent) error {
	log.Printf("name: %v", name)
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to retrieve IP address; StatusCode: %d", resp.StatusCode)
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String("ip-address"),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   bytes.NewReader(ip),
	}

	_, err = StoreIpAddress(context.TODO(), client, input)
	return err

}

func main() {
	lambda.Start(handler)
}
