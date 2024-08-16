package utils

import "os"

// ReadFile 读取文件
func ReadFile(file string) ([]byte, error) {
	if f, err := os.Open(file); err != nil {
		return nil, err
	} else {
		content := make([]byte, 4096)
		if n, err := f.Read(content); err != nil {
			return nil, err
		} else {
			return content[:n], err
		}
	}
}

// ReadFileString 读取文件
func ReadFileString(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
