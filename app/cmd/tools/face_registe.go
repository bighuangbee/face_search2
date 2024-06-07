package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var faceServer1 = "http://localhost:6002"

type RegisteRespone struct {
	Code    int     `json:"code"`
	Data    Registe `json:"data"`
	Message string  `json:"message"`
	Reason  string  `json:"reason"`
}

type Registe struct {
	NewFaceNum        int `json:"newFaceNum"`
	RegistedFailedNum int `json:"registedFailedNum"`
	RegistedSuccNum   int `json:"registedSuccNum"`
}

func main() {

	var errPublic error = nil

	url := faceServer1 + "/face/registe/path"
	data := `{}`
	payload := bytes.NewBufferString(data)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {
		registeRespone := RegisteRespone{}
		if err := json.Unmarshal(respBody, &registeRespone); err != nil {
			fmt.Println("解析返回数据失败:", err, string(respBody))
			errPublic = err
		}

		fmt.Println("当前新增注册的人脸数:", registeRespone.Data.NewFaceNum,
			", 已注册失败的人脸数:", registeRespone.Data.RegistedFailedNum,
			", 已注册成功的人脸数:", registeRespone.Data.RegistedSuccNum)

	} else {
		errPublic = errors.New("StatusCode " + strconv.Itoa(resp.StatusCode))
	}

	if errPublic != nil {
		for {
		}
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		status, err := getRegisteStatus()
		if err != nil {
			fmt.Println("检查人脸注册状态错误", err)
			continue
		}

		if status {
			fmt.Println("检查人脸注册状态", "正在注册...")
		} else {
			return
		}
	}

}

type RegisteStatusRespone struct {
	Code    int           `json:"code"`
	Data    RegisteStatus `json:"data"`
	Message string        `json:"message"`
	Reason  string        `json:"reason"`
}

type RegisteStatus struct {
	Registering bool `json:"registering"`
}

func getRegisteStatus() (bool, error) {
	url := faceServer1 + "/face/registe/status"

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusOK {
		registeRespone := RegisteStatusRespone{}
		if err := json.Unmarshal(respBody, &registeRespone); err != nil {
			fmt.Println("解析返回数据失败:", err, string(respBody))
			return false, err
		}

		return registeRespone.Data.Registering, nil
	}

	return false, errors.New("getRegisteStatus error")
}
