package s3

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	downloadThreads = 50
)

type downloadJob struct {
	key  string
	path string
}

func (s *S3) DownloadFile(ctx context.Context, key, dir string) error {
	if keyIsDir(key) {
		return errors.New("cant download folder")
	}

	path := getDownloadPath(key, dir, "")

	return s.downloadKey(ctx, key, path)
}

func (s *S3) DownloadFolder(ctx context.Context, prefix, dir string) error {
	if !keyIsDir(prefix) {
		return errors.New("cant download file")
	}

	keys, err := s.listForDownload(ctx, prefix)
	if err != nil {
		return err
	}

	downloadJobs := make(chan downloadJob)

	for i := 0; i < downloadThreads; i++ {
		go s.downloadWorker(ctx, downloadJobs)
	}

	scheduledKeys := 0

	for _, key := range keys {
		scheduledKeys++
		percent := 100 * scheduledKeys / len(keys)
		runtime.EventsEmit(ctx, "download-process", percent)

		downloadJobs <- downloadJob{
			key:  key,
			path: getDownloadPath(key, dir, prefix),
		}
	}

	return nil
}

func (s *S3) listForDownload(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	var after string

	for {
		i := s3.ListObjectsV2Input{
			Bucket:  aws.String(s.bucket),
			Prefix:  aws.String(prefix),
			MaxKeys: aws.Int64(maxKeys),
		}

		if after != "" {
			i.StartAfter = aws.String(after)
		}

		res, err := s.client.ListObjectsV2WithContext(ctx, &i)
		if err != nil {
			return keys, err
		}

		for _, object := range res.Contents {
			keys = append(keys, *object.Key)
			after = *object.Key
		}

		if len(res.Contents) < maxKeys {
			break
		}
	}

	return keys, nil
}

func (s *S3) downloadWorker(ctx context.Context, jobs <-chan downloadJob) {
	for j := range jobs {
		_ = s.downloadKey(ctx, j.key, j.path)
	}
}

func (s *S3) downloadKey(ctx context.Context, key, path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	res, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	defer func() { _ = outFile.Close() }()
	_, err = io.Copy(outFile, res.Body)

	return nil
}
