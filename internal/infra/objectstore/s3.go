package objectstore

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Options struct {
	Region    string
	Endpoint  string // Optional S3-compatible endpoint, e.g. MinIO.
	PathStyle bool
}

// OpenS3 loads the standard AWS credential chain and creates a reusable client.
// It does not issue a network request; call CheckBucket for a required bucket.
func OpenS3(ctx context.Context, opts S3Options) (*s3.Client, error) {
	if opts.Region == "" {
		return nil, errors.New("S3 region is required")
	}
	if opts.Endpoint != "" {
		endpoint, err := url.Parse(opts.Endpoint)
		if err != nil || endpoint.Host == "" || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" || endpoint.Scheme != "https" && endpoint.Scheme != "http" {
			return nil, errors.New("S3 endpoint must be an absolute HTTP(S) URL")
		}
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(opts.Region))
	if err != nil {
		return nil, fmt.Errorf("load S3 configuration: %w", err)
	}
	return s3.NewFromConfig(cfg, func(client *s3.Options) {
		if opts.Endpoint != "" {
			client.BaseEndpoint = aws.String(opts.Endpoint)
		}
		client.UsePathStyle = opts.PathStyle
	}), nil
}

// CheckBucket verifies access to one required bucket. HeadBucket requires the
// relevant bucket permission; use it for readiness only if access is required.
func CheckBucket(ctx context.Context, client *s3.Client, bucket string) error {
	if client == nil || bucket == "" {
		return errors.New("S3 client and bucket are required")
	}
	if _, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return fmt.Errorf("check S3 bucket: %w", err)
	}
	return nil
}
