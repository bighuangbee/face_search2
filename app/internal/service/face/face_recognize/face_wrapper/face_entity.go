package face_wrapper

/*

#cgo LDFLAGS: -L../../../../../../libs/sdk/lib/ -lhiar_cluster
#cgo CFLAGS: -I ./
#include "./interface_face_recognizer.h"

#include <stdbool.h>
#include <stdlib.h>
*/
import "C"
import (
	"path"
	"strings"
	"unsafe"
)

type Personnal struct {
	Type string //人脸类型
	ID   string //身份证id
	Name string //姓名
}

type FaceEntity struct {
	RegFilename string //人脸注册文件名 aa/bb.jpg
	Personnal   *Personnal

	Match float32 //比对匹配值
}

const HIAR_FACE_FEATURE_LEN = 512

//================图像数据======================

// 人脸输入类型
const (
	IMAGE_TYPE_JPG = "jpg"
	IMAGE_TYPE_PNG = "png"
)

var FaceDataType = map[string]ImageType{
	IMAGE_TYPE_JPG: ImageTypeJPG,
	IMAGE_TYPE_PNG: ImageTypePNG,
}

type ImageType int

const (
	ImageTypeJPG ImageType = 1
	ImageTypePNG ImageType = 2
)

type Image struct {
	DataType ImageType
	Size     int
	Data     []byte
	Width    int
	Height   int
}

func GetImageType(filePath string) ImageType {
	if t, ok := FaceDataType[strings.Trim(path.Ext(filePath), ".")]; ok {
		return ImageType(t)
	}
	return ImageTypeJPG
}

/**
 * @Description: 创建图像 cgo内存
 * @param data
 * @param width
 * @param height
 * @param dataType
 * @return *C.ImageData 图像数据类型
 */
func NewC_ImageData(image *Image) *C.ImageData {

	var imageData C.ImageData
	//todo free

	imageData.data = (*C.uchar)((unsafe.Pointer)(&image.Data[0]))
	imageData.data_len = C.int(len(image.Data))
	imageData.width = C.int(image.Width)
	imageData.height = C.int(image.Height)
	imageData.data_type = C.enum_ImageDataType(image.DataType)
	return &imageData

}
