package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type AWSProvider struct {
	s3Client         *s3.Client
	awsConfiguration aws.Config
}

// NewAWSProvider creates a new AWS provider instance
func NewAWSProvider(region, accessKeyID, secretAccessKey string) (AWSProvider, error) {
	// Create custom credentials
	awsCredentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		),
	)

	// Load AWS configuration with custom credentials
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(awsCredentials),
	)
	if err != nil {
		return AWSProvider{}, err
	}

	return AWSProvider{
		s3Client:         s3.NewFromConfig(cfg),
		awsConfiguration: cfg,
	}, nil
}

// S3 Implementation
func (p *AWSProvider) UploadToS3(ctx context.Context, bucketName, region string, key string, data []byte) (*string, error) {
	_, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		logrus.Tracef("file upload failed %v", err)
		err = fmt.Errorf("could not upload file %v", err)
		return nil, err
	}
	url := fmt.Sprintf(`https://%v.s3.%v.amazonaws.com/%v`, bucketName, region, key)
	return &url, nil
}

func (p *AWSProvider) GetS3Object(ctx context.Context, bucketName, key string) ([]byte, error) {
	result, err := p.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (p *AWSProvider) DeleteS3Object(ctx context.Context, bucketName, key string) error {
	_, err := p.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	return err
}
