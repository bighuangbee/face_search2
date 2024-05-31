package main

import (
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"io/ioutil"
)

var basePath = "../../../../../"

func main() {
	err := face_wrapper.Init("/root/face_search/libs/models/", "./hiarClusterLog.txt")
	if err != nil {
		panic(err)
	}

	regNum, err := face_wrapper.Registe(basePath + "/libs/data/gallery")
	if err != nil {
		fmt.Println("face_wrapper.Registe, error", err)
		return
	}

	fmt.Println("face_wrapper.Registe, regNum", regNum)

	imageFile, err := ioutil.ReadFile(basePath + "libs/data/query.jpg")
	if err != nil {
		fmt.Println("ReadFile", err)
		return
	}

	targetFace := face_wrapper.Image{
		DataType: face_wrapper.GetImageType(basePath + "libs/data/query.jpg"),
		Size:     len(imageFile),
		Data:     imageFile,
	}

	results := face_wrapper.Search(&targetFace)
	if len(results) > 0 {
		for key, result := range results {
			fmt.Println("face_wrapper.Search result:", key+1, result.RegFilename, result.Match)
		}
	} else {
		fmt.Println("搜索不到结果")
	}

	face_wrapper.UnRegisteAll()

	results2 := face_wrapper.Search(&targetFace)
	if len(results2) > 0 {
		for key, result := range results2 {
			fmt.Println("face_wrapper.Search result:", key+1, result.RegFilename, result.Match)
		}
	} else {
		fmt.Println("搜索不到结果")
	}
}
