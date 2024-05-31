package util

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFilesWithExtensions 读取指定目录下的文件，并过滤出指定扩展名的文件
func GetFilesWithExtensions(dir string, extensions []string) ([]string, error) {
	var files []string

	// Walk 函数遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是文件，并且扩展名是否在指定的扩展名列表中
		if !info.IsDir() && HasValidExtension(path, extensions) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// HasValidExtension 检查文件是否具有有效的扩展名
func HasValidExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.EqualFold(filepath.Ext(file), ext) {
			return true
		}
	}
	return false
}

func GetFileName(filePath string) string {
	fileName := filepath.Base(filePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
