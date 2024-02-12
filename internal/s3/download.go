package s3

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type downloadJob struct {
	key  string
	path string
}

func (s *S3) Download(ctx context.Context, prefix, downloadDir string) error {
	prefix = strings.TrimLeft(prefix, "/")
	downloadDir = strings.TrimRight(downloadDir, "/")

	pParts := strings.Split(prefix, "/")
	isDir := pParts[len(pParts)-1] == ""
	if !isDir {
		return s.downloadKey(ctx, prefix, downloadDir+"/"+pParts[len(pParts)-1])
	}

	dirName := ""
	if len(pParts) > 1 {
		dirName = pParts[len(pParts)-2]
	}

	baseDownloadPath := downloadDir + "/"
	if dirName != "" {
		baseDownloadPath += dirName + "/"
	}

	l, err := s.List(ctx, prefix)
	if err != nil {
		return err
	}

	downloadJobs := make(chan downloadJob)

	for i := 0; i < downloadThreads; i++ {
		go s.downloadWorker(ctx, downloadJobs)
	}

	scheduledKeys := 0

	for _, key := range l.Files {
		keyParts := strings.Split(key, "/")
		if keyParts[len(keyParts)-1] == "" {
			continue
		}

		scheduledKeys++
		percent := 100 * scheduledKeys / len(l.Files)
		runtime.EventsEmit(ctx, "download-process", percent)

		downloadJobs <- downloadJob{
			key:  key,
			path: baseDownloadPath + strings.Replace(key, prefix, "", 1),
		}
	}

	return nil
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
