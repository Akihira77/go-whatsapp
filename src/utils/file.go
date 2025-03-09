package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func GetUploadDestination() (string, error) {
	path, err := os.Getwd()

	return fmt.Sprintf("%s/storage/uploads/", path), err
}

func ReadFile(fileHeader *multipart.FileHeader) (*bytes.Buffer, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)

	return buf, err
}
