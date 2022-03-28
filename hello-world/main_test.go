package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

type mockPutObjectAPI func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)

func (m mockPutObjectAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func Test_storeIpAddress(t *testing.T) {

	os.Setenv("BUCKET_NAME", "test-bucket")

	type args struct {
		api S3PutObjectAPI
		ip  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Store at ip-address",
			args: args{
				api: mockPutObjectAPI(func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					t.Helper()
					assert.NotNil(t, params.Bucket, "expect bucket to not be nil")
					assert.Equal(t, os.Getenv("BUCKET_NAME"), *params.Bucket)
					assert.NotNil(t, params.Key, "expect key to not be nil")
					assert.Equal(t, "ip-address", *params.Key)
					assert.EqualValues(t, "public-read", params.ACL)
					// t.Logf("params: %v", params.ACL)

					return &s3.PutObjectOutput{}, nil
				}),
				ip: []byte("1.1.1.1"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeIpAddress(tt.args.api, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("storeIpAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
