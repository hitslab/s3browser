package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	maxKeys = 1000
)

type List struct {
	Folders []string `json:"folders"`
	Files   []string `json:"files"`
}

func (s *S3) List(ctx context.Context, prefix string) (List, error) {
	var l List
	var after string

	for {
		i := s3.ListObjectsV2Input{
			Bucket:    aws.String(s.bucket),
			Prefix:    aws.String(prefix),
			Delimiter: aws.String("/"),
			MaxKeys:   aws.Int64(maxKeys),
		}

		if after != "" {
			i.StartAfter = aws.String(after)
		}

		res, err := s.client.ListObjectsV2WithContext(ctx, &i)
		if err != nil {
			return l, err
		}

		for _, object := range res.CommonPrefixes {
			l.Folders = append(l.Folders, *object.Prefix)
		}

		for _, object := range res.Contents {
			l.Files = append(l.Files, *object.Key)
			after = *object.Key
		}

		if len(res.Contents) < maxKeys {
			break
		}
	}

	return l, nil
}
