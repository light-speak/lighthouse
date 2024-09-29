package file

import (
	"os"
)

type FileOutput struct {
	file *os.File
}

func NewFileOutput(path string) (*FileOutput, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileOutput{file: f}, nil
}

func (fo *FileOutput) Write(p []byte) (n int, err error) {
	return fo.file.Write(p)
}