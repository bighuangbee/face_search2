package face

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fileInfoList := LoadFileInfo()

	startTime := "2024-06-21" + " 17:47"
	endTime := "2024-06-21" + " 17:57"

	results, _ := GetRangeFile(fileInfoList, startTime, endTime)
	for i, result := range results {
		fmt.Println("result", i+1, *result)
	}
}
