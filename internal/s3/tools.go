package s3

import (
	"fmt"
	"path/filepath"
	"strings"
)

func keyIsDir(prefix string) bool {
	prefix = strings.TrimLeft(prefix, "/")
	pParts := strings.Split(prefix, "/")
	return pParts[len(pParts)-1] == ""
}

func getDownloadPath(key, dir string, p string) string {
	dir = strings.TrimRight(dir, "/") + "/"

	if p == "" {
		return fmt.Sprintf("%s%s", dir, filepath.Base(key))
	}

	p = strings.TrimRight(p, "/") + "/"

	folder := ""
	pParts := strings.Split(p, "/")
	if len(pParts) >= 2 && pParts[len(pParts)-2] != "" {
		folder = pParts[len(pParts)-2] + "/"
	}

	fileName := strings.Replace(key, p, "", 1)

	return fmt.Sprintf("%s%s%s", dir, folder, fileName)
}
