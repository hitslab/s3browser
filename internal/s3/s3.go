package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hitslab/s3browser/internal/config"
)

type S3 struct {
	client *s3.S3
	bucket string
}

func (s *S3) Connect(ctx context.Context, st config.S3Settings) error {
	ses, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(st.AccessKey, st.SecretKey, ""),
		Endpoint:         aws.String(st.BaseUrl),
		Region:           aws.String(st.Region),
		DisableSSL:       aws.Bool(st.DisableSSL),
		S3ForcePathStyle: aws.Bool(st.PathStyle),
	})

	if err != nil {
		return err
	}

	s.client = s3.New(ses)

	_, err = s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(st.Bucket),
		MaxKeys: aws.Int64(5),
	})

	if err != nil {
		return err
	}

	s.bucket = st.Bucket

	return nil
}
