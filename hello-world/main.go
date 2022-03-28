package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// we use S3PutObjectAPI because we can mock it in our tests
func storeIpAddress(api S3PutObjectAPI, ip []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")), // test env is used
		Key:    aws.String("ip-address"),             // ensure object name
		ACL:    types.ObjectCannedACLPublicRead,      // ensure ACL is "public-read"
		Body:   bytes.NewReader(ip),
	}
	_, err := api.PutObject(context.TODO(), input)
	return err
}

func whatIsMyIp(url string) (ip []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return ip, err
	}

	// only "business logic", test an error thrown if the response is not 200 - aka error handling
	if resp.StatusCode != 200 {
		return ip, fmt.Errorf("failed to retrieve IP address; StatusCode: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func handler(ctx context.Context) error {

	ip, err := whatIsMyIp("https://checkip.amazonaws.com")
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	return storeIpAddress(client, ip)
}

func main() {
	lambda.Start(handler)
}
