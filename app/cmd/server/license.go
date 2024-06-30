package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TencentTimeAPIResponse struct {
	Success string `json:"success"`
	Result  struct {
		Datetime1 string `json:"datetime_1"`
	} `json:"result"`
}

func getCurrentTimeFromTencent() (time.Time, error) {
	resp, err := http.Get("http://api.k780.com/?app=life.time&appkey=10003&sign=b59bc3ef6191eb9f747dd4e83c99f2a4&format=json")
	if err != nil {
		return time.Time{}, fmt.Errorf("获取时间时出错: %v", err)
	}
	defer resp.Body.Close()

	var result TencentTimeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return time.Time{}, fmt.Errorf("解析时间响应时出错: %v", err)
	}

	currentTime, err := time.Parse("2006-01-02 15:04:05", result.Result.Datetime1)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间格式时出错: %v", err)
	}

	return currentTime, nil
}

func isAfterSpecificDate() (bool, error) {
	specificDateStr := "2024-07-30 23:59:59"

	specificDate, err := time.Parse(time.DateTime, specificDateStr)
	if err != nil {
		return true, fmt.Errorf("解析特定日期时出错: %v", err)
	}

	now, err := getCurrentTimeFromTencent()
	if err != nil {
		fmt.Println("000now ", now)
		return true, fmt.Errorf("获取当前时间时出错: %v", err)
	}

	return now.After(specificDate), nil
}
