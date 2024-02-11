package s3

import (
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

type List struct {
	Folders []string `json:"folders"`
	Files   []string `json:"files"`
}

const (
	batchSize = 1000
	maxFiles  = 100000
)

func (s *S3) Connect(st config.S3Settings) error {
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

	_, err = s.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(st.Bucket),
		MaxKeys: aws.Int64(5),
	})

	if err != nil {
		return err
	}

	s.bucket = st.Bucket

	return nil
}

func (s *S3) List(prefix string) (List, error) {
	var l List
	var listObjects *s3.ListObjectsV2Output
	var err error

	delimiter := "/"
	after := ""

	for {
		i := s3.ListObjectsV2Input{
			Bucket:    aws.String(s.bucket),
			Prefix:    aws.String(prefix),
			Delimiter: aws.String(delimiter),
			MaxKeys:   aws.Int64(batchSize),
		}

		if after != "" {
			i.StartAfter = aws.String(after)
		}

		listObjects, err = s.client.ListObjectsV2(&i)
		if err != nil {
			return List{}, err
		}

		for _, object := range listObjects.CommonPrefixes {
			l.Folders = append(l.Folders, *object.Prefix)
		}

		for _, object := range listObjects.Contents {
			l.Files = append(l.Files, *object.Key)
			after = *object.Key
		}

		if len(l.Files) >= maxFiles || len(listObjects.Contents) < batchSize {
			break
		}
	}

	return l, nil
}
