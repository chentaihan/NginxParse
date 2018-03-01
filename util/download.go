package util

import (
	"path/filepath"
	"os"
	"io"
	"archive/zip"
)

func Downalod(url ,filePath string) bool{
	buffer := HttpGet(url)
	if buffer != nil {
		err := WriteFile(filePath, buffer, 0755)
		if err == nil {
			return true
		}
	}
	return false
}

func UnZip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		copyFile(target, file)
	}

	return nil
}

func copyFile(target string, file *zip.File) bool {
	path := filepath.Join(target, file.Name)
	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
		return true
	}

	fileReader, err := file.Open()
	if err != nil {
		return false
	}
	defer fileReader.Close()

	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return false
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return false
	}
	return true
}
