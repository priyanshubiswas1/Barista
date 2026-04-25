package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type GarageUploader struct {
	bucket          string
	publicURLPrefix string
	uploader        *manager.Uploader
}

type GarageConfig struct {
	Endpoint        string
	Region          string
	AccessKey       string
	SecretKey       string
	Bucket          string
	PublicURLPrefix string
}

func NewGarageUploader(ctx context.Context, cfg GarageConfig) (*GarageUploader, error) {
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithHTTPClient(nil),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if service == s3.ServiceID && cfg.Endpoint != "" {
					return aws.Endpoint{URL: cfg.Endpoint, HostnameImmutable: true}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &GarageUploader{
		bucket:          cfg.Bucket,
		publicURLPrefix: cfg.PublicURLPrefix,
		uploader:        manager.NewUploader(s3Client),
	}, nil
}

func (g *GarageUploader) UploadBytes(ctx context.Context, key string, contentType string, data []byte) (string, error) {
	_, err := g.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(g.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// If a public prefix is configured, return a URL; otherwise return "s3://bucket/key".
	if strings.TrimSpace(g.publicURLPrefix) != "" {
		prefix := strings.TrimRight(g.publicURLPrefix, "/")
		return fmt.Sprintf("%s/%s", prefix, key), nil
	}
	return fmt.Sprintf("s3://%s/%s", g.bucket, key), nil
}

func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

