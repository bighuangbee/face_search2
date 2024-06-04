package main

import (
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"io/ioutil"
)

var basePath = "../../../../../"

func main() {
	err := face_wrapper.Init("/root/face_search/libs/models/", "./hiarClusterLog.txt")
	if err != nil {
		panic(err)
	}

	files, err := util.GetFilesWithExtensions(basePath+"/libs/data/gallery", face_wrapper.PictureExt)
	if err != nil {
		fmt.Println("GetFilesWithExtensions, error", err)
		return
	}

	for _, filename := range files {
		fmt.Println("filename", filename)
		imageFile, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("ReadFile", err)
			return
		}
		if err := face_wrapper.RegisteSingle(&face_wrapper.Image{
			DataType: face_wrapper.GetImageType(filename),
			Size:     len(imageFile),
			Data:     imageFile,
		}, filename); err != nil {
			fmt.Println("face_wrapper.Registe, error", err)
			return
		}

		fmt.Println("face_wrapper.RegisteSingle", err)
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

	//face_wrapper.Search2(&targetFace)

	results := face_wrapper.Search(&targetFace)
	if len(results) > 0 {
		for key, result := range results {
			fmt.Println("【Search】face_wrapper.Search result:", key+1, result.RegFilename, result.Match)
		}
	} else {
		fmt.Println("搜索不到结果")
	}

	//face_wrapper.UnRegisteAll()
	//
	//results2 := face_wrapper.Search(&targetFace)
	//if len(results2) > 0 {
	//	for key, result := range results2 {
	//		fmt.Println("face_wrapper.Search result:", key+1, result.RegFilename, result.Match)
	//	}
	//} else {
	//	fmt.Println("搜索不到结果")
	//}
}
