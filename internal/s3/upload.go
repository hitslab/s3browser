package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (s *S3) UploadFiles(ctx context.Context, prefix string, files []string) error {
	for _, f := range files {
		key := strings.TrimRight(prefix, "/") + "/" + filepath.Base(f)
		err := s.uploadFile(ctx, key, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *S3) UploadFolder(ctx context.Context, prefix string, dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		key := strings.TrimRight(prefix, "/") +
			"/" +
			filepath.Base(dir) +
			strings.Replace(path, dir, "", 1)

		return s.uploadFile(ctx, key, path)
	})

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (s *S3) uploadFile(ctx context.Context, key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}
