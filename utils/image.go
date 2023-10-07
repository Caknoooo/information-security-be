package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

const (
	PATH = "storage"
)

func UploadFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func UploadFileSuccess(file *multipart.FileHeader, path string) error {
	parts := strings.Split(path, "/")

	fileId := parts[2]
	directoryPath := fmt.Sprintf("%s/%s/%s", PATH, parts[0], parts[1])

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, 0777); err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s", directoryPath, fileId)

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	fileData, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, fileData, 0666)
	if err != nil {
		return err
	}

	return nil
}
