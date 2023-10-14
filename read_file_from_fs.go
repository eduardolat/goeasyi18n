package goeasyi18n

import (
	"io"
	"io/fs"
)

func readFileFromFS(filesystem fs.FS, file string) ([]byte, error) {
	fileData, err := filesystem.Open(file)
	if err != nil {
		return nil, err
	}
	defer fileData.Close()

	return io.ReadAll(fileData)
}
