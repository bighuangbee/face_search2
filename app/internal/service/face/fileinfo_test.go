package face

import (
	"fmt"
	"testing"
	"time"
)

// 模拟文件信息结构
//type FileInfo struct {
//	Birthtime time.Time
//}

// LoadFileInfo 模拟加载文件信息的函数
//func LoadFileInfo() map[string]FileInfo {
//	return map[string]FileInfo{
//		"file1.txt": {Birthtime: time.Date(2024, 6, 21, 17, 47, 0, 0, time.Local)},
//		"file2.txt": {Birthtime: time.Date(2024, 6, 21, 17, 48, 0, 0, time.Local)},
//		"file3.txt": {Birthtime: time.Date(2024, 6, 21, 17, 49, 0, 0, time.Local)},
//	}
//}

func TestName(t *testing.T) {
	fileInfoList := LoadFileInfo()

	inputTimeStr := "2024-06-21" + " 17:47"

	timeRange := time.Minute * 1
	results, _ := GetRangeFile(fileInfoList, inputTimeStr, timeRange)
	for i, result := range results {
		fmt.Println("result", i+1, *result)
	}
}
