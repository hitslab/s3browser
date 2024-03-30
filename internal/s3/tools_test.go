package s3

import (
	"fmt"
	"strconv"
	"testing"
)

func TestKeyIsDir(t *testing.T) {
	type testCase struct {
		key       string
		wantIsDir bool
	}

	tests := []testCase{
		{
			"/some/key", false,
		},
		{
			"/", true,
		},
		{
			"", true,
		},
		{
			"/some/prefix/", true,
		},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			isDir := keyIsDir(test.key)
			if isDir != test.wantIsDir {
				t.Fatalf(`result is incorrect, got: "%s", want "%s"`, strconv.FormatBool(isDir), strconv.FormatBool(test.wantIsDir))
			}
		})
	}
}

func TestGetDownloadPath(t *testing.T) {
	type testCase struct {
		key      string
		dir      string
		prefix   string
		wantPath string
	}

	tests := []testCase{
		{
			"/file.txt", "/download", "",
			"/download/file.txt",
		},
		{
			"/some/prefix/file.txt", "/download", "",
			"/download/file.txt",
		},
		{
			"/some/prefix/file.txt", "/download/", "",
			"/download/file.txt",
		},
		{
			"/some/prefix/file.txt", "/download/some/folder", "",
			"/download/some/folder/file.txt",
		},
		{
			"/some/prefix/file.txt", "", "",
			"/file.txt",
		},
		{
			"/some/prefix/file.txt", "/", "",
			"/file.txt",
		},
		{
			"/some/long/prefix/file.txt", "/download", "/some",
			"/download/some/long/prefix/file.txt",
		},
		{
			"/some/long/prefix/file.txt", "/download/", "/some/",
			"/download/some/long/prefix/file.txt",
		},
		{
			"/some/long/prefix/file.txt", "/download/folder", "/some/",
			"/download/folder/some/long/prefix/file.txt",
		},
		{
			"/some/long/prefix/file.txt", "/download/folder", "/",
			"/download/folder/some/long/prefix/file.txt",
		},
		{
			"/some/long/prefix/file.txt", "/download/folder", "/some/long/prefix/",
			"/download/folder/prefix/file.txt",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			path := getDownloadPath(test.key, test.dir, test.prefix)
			if path != test.wantPath {
				t.Fatalf(`result path is incorrect, got: "%s", want "%s"`, path, test.wantPath)
			}
		})
	}
}
