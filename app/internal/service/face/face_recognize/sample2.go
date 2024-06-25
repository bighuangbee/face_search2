package main

import (
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"io/ioutil"
)

var basePath = "../../../../../"

func main() {
	err := face_wrapper.Init("/root/face_search/libs/models/", 0.8, "./")
	if err != nil {
		panic(err)
	}

	files, err := util.GetFilesWithExtensions(basePath+"/libs/data/gallery", face_wrapper.PictureExt)
	if err != nil {
		fmt.Println("GetFilesWithExtensions, error", err)
		return
	}

	var regInfo []*face_wrapper.RegisteInfo
	for _, filename := range files {
		regInfo = append(regInfo, &face_wrapper.RegisteInfo{
			Filename: filename,
		})
	}

	if err := face_wrapper.Registe(regInfo); err != nil {
		fmt.Println("face_wrapper.Registe", err)
		return
	}

	for i, info := range regInfo {
		fmt.Println("注册结果", i+1, info.Time, info.Ok, info.Filename)
	}

	targetFile := basePath + "libs/data/query.jpg"
	//targetFile := basePath + "libs/data/gallery/DSC08060.JPG"
	imageFile, err := ioutil.ReadFile(targetFile)
	if err != nil {
		fmt.Println("ReadFile", err)
		return
	}

	targetFace := face_wrapper.Image{
		DataType: face_wrapper.GetImageType(targetFile),
		Size:     len(imageFile),
		Data:     imageFile,
	}

	results := face_wrapper.Search(&targetFace)
	if len(results) > 0 {
		for key, result := range results {
			fmt.Println("【Search】人脸检索结果:", key+1, result.RegFilename, result.Match)
		}
	} else {
		fmt.Println("搜索不到结果")
	}

	fmt.Println("注销全部人脸")

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
