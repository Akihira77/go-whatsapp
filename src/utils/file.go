package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func GetUploadDestination() (string, error) {
	path, err := os.Getwd()

	return fmt.Sprintf("%s/storage/uploads/", path), err
}

func SaveUploadedFile(ctx context.Context, file *multipart.FileHeader, dst string) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return nil, err
	}

	out, err := os.Create(dst + file.Filename)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	buf, err := ReadFile(src)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ReadFile(file multipart.File) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)

	return buf, err
}
