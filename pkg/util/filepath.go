package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetFilesWithExtensions 读取指定目录下的文件，并过滤出指定扩展名的文件
func ReadFilesWithExtensions(dir string, extensions []string) ([]string, error) {
	fileList := []string{}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if !info.IsDir() && HasValidExtension(file.Name(), extensions) {
			fileList = append(fileList, file.Name())
		}

	}
	return fileList, nil
}

// GetFilesWithExtensions 读取指定目录下的文件，并过滤出指定扩展名的文件
func GetFilesWithExtensions(dir string, extensions []string) ([]string, error) {
	var files []string

	// Walk 函数遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 转换路径分隔符为正斜杠
		path = filepath.ToSlash(path)

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

func DecodeGBKToUTF8(input []byte) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, input)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

// IsUTF8 checks if a byte slice is valid UTF-8 encoded
func IsUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// DecodeToUTF8 tries to decode a byte slice using the provided decoder
func DecodeToUTF8(data []byte, decoder transform.Transformer) (string, error) {
	utf8Bytes, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

// detectAndDecode tries to detect if the input is in GBK or GB2312, and decodes to UTF-8
func DetectAndDecode(input []byte) (string, error) {
	// First, check if it's already UTF-8
	if IsUTF8(input) {
		return string(input), nil
	}

	// Try decoding with GBK
	decodedStr, err := DecodeToUTF8(input, simplifiedchinese.GBK.NewDecoder())
	if err == nil {
		return decodedStr, nil
	}

	// Try decoding with HZGB2312
	decodedStr, err = DecodeToUTF8(input, simplifiedchinese.HZGB2312.NewDecoder())
	if err == nil {
		return decodedStr, nil
	}

	return "", fmt.Errorf("failed to decode input: %v", input)
}

// 创建或打开当前日期的文件，并写入内容
func CreateOrOpenFile(filename, content string) error {
	// 判断文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close() // 确保文件最后被关闭
	}

	// 打开文件进行写入
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入内容
	_, err = file.WriteString(content + "\n")
	return err
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
