package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

var storagepath string

func init() {
	storagepath = os.Getenv("FILE_STORAGE_PATH")
}

//GetFilename function
func GetFilename(uid, prefix, ext string) (filename, to string) {
	filename = fmt.Sprintf("%s-%d%s", uid, time.Now().Unix(), ext)
	to = fmt.Sprintf("%s/%s/%s", storagepath, prefix, filename)
	return filename, to
}

// SaveFile function
func SaveFile(file multipart.File, headers *multipart.FileHeader, uid, prefix string) (string, error) {
	ext := filepath.Ext(headers.Filename)
	filename, to := GetFilename(uid, prefix, ext)
	f, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return filename, nil
}
