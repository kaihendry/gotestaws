package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

type mockPutObjectAPI func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)

func (m mockPutObjectAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestStoreIpAddress(t *testing.T) {

	s3Bucket := "test_bucket"

	type args struct {
		c     context.Context
		api   mockPutObjectAPI
		input *s3.PutObjectInput
	}
	tests := []struct {
		name    string
		args    args
		want    *s3.PutObjectOutput
		wantErr bool
	}{
		{
			name: "First",
			args: args{
				c: nil,
				api: mockPutObjectAPI(func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					t.Helper()
					assert.NotNil(t, params.Bucket, "expect bucket to not be nil")
					assert.Equal(t, s3Bucket, *params.Bucket)
					assert.NotNil(t, params.Key, "expect key to not be nil")
					assert.Equal(t, "ip-address", *params.Key)

					return &s3.PutObjectOutput{}, nil
				}),
				input: &s3.PutObjectInput{
					Bucket: aws.String(s3Bucket),
					Key:    aws.String("ip-address"),
					ACL:    types.ObjectCannedACLPublicRead,
				},
			},
			want:    &s3.PutObjectOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreIpAddress(tt.args.c, tt.args.api, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreIpAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreIpAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
