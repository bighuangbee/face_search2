package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bighuangbee/face_search2/pkg/util"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var PictureExt = []string{".png", ".jpg", ".jpeg"}
var faceServer = "http://localhost:6002"
var registePath = "C:\\hiar_face\\registe_path"

type SearchRespone struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type Results struct {
	Filename string  `json:"filename"`
	Match    float64 `json:"match"`
}

type Data struct {
	Results []Results `json:"results"`
}

func main() {

	files, err := util.ReadFilesWithExtensions("./", PictureExt)
	if err != nil {
		fmt.Println("GetFilesWithExtensions, error", err)
		return
	}

	var errPublic error = nil

	for i, filename := range files {
		fmt.Println("搜索人脸 ", i+1, filename)

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("打开文件失败:", err)
			continue
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			fmt.Println("创建form file失败:", err)
			continue
		}

		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println("写入文件内容失败:", err)
			continue
		}

		err = writer.Close()
		if err != nil {
			fmt.Println("关闭writer失败:", err)
			continue
		}

		t1 := time.Now()
		// 创建HTTP请求
		req, err := http.NewRequest("POST", faceServer+"/face/search", body)
		if err != nil {
			fmt.Println("创建HTTP请求失败:", err)
			continue
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("发送HTTP请求失败:", err)
			continue
		}
		defer resp.Body.Close()

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应失败:", err)
			continue
		}

		resultPath := "result_" + util.GetFileName(filename)

		if resp.StatusCode == http.StatusOK {
			searchRespone := SearchRespone{}
			if err := json.Unmarshal(respBody, &searchRespone); err != nil {
				fmt.Println("解析返回数据失败:", err, string(respBody))
				errPublic = err
			}

			if _, err := os.Stat(resultPath); !os.IsNotExist(err) {
				err := os.RemoveAll(resultPath)
				if err != nil {
					fmt.Printf("Error removing directory: %s\n", err)
				} else {
					fmt.Printf("Directory '%s' has been deleted.\n", resultPath)
				}
			}

			fmt.Println(len(searchRespone.Data.Results), searchRespone.Data.Results)
			if len(searchRespone.Data.Results) > 0 {
				//createResultPath(resultPath)
				fmt.Println("创建目录", os.Mkdir(resultPath, 0644), resultPath)

				fmt.Println("复制目标图片", copyFile(filename, resultPath+"\\targer_"+filename), resultPath+"\\targer_"+filename)

				for _, result := range searchRespone.Data.Results {
					fmt.Println("result    === ", result.Filename, result.Match)
					src := registePath + "\\" + result.Filename
					dst := resultPath + "\\" + fmt.Sprintf("%s_%f%s", util.GetFileName(filename), result.Match, filepath.Ext(result.Filename))
					err := copyFile(src, dst)
					fmt.Println("复制人脸搜索结果", err, "src:", src, "dst:", dst)
				}
			}

		} else {
			errPublic = errors.New("StatusCode " + strconv.Itoa(resp.StatusCode))
		}

		//耗时
		os.Create(resultPath + "\\" + fmt.Sprintf("druation_%dms.txt", time.Since(t1).Milliseconds()))
	}

	if errPublic != nil {
		fmt.Println("error", errPublic)
	} else {
		fmt.Println("搜索完成")
	}
	for {
	}
}

func createResultPath(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Error removing directory: %s\n", err)
		} else {
			fmt.Printf("Directory '%s' has been deleted.\n", path)
		}

		fmt.Println("创建目录", os.Mkdir(path, 0644), path)
	} else {
		fmt.Printf("Directory '%s' does not exist.\n", path)
	}
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return
}
