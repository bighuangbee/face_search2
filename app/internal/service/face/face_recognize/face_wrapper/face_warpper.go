package face_wrapper

/*
#cgo LDFLAGS: -L../libs/sdk/lib/ -lhiar_cluster
#cgo CFLAGS: -I ./
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "./interface_face_recognizer.h"


int modelPath(const char* model_path){
	printf("modelPath   Loading model from: %s\n", model_path);
}

void printImageInfo(ImageInfo* info, int count) {
    for (int i = 0; i < count; i++) {
        printf("C.ImageInfo %d: filename = %s, similarity = %f\n", i, info[i].filename, info[i].similarity);
    }
}

*/
import "C"
import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"unsafe"
)

const FACE_MAX_RESULT = 20
const OK = C.int(1)

var PictureExt = []string{".png", ".jpg", ".jpeg"}

func Init(modelPath string, logFilename string) error {
	ret := C.hiarClusterInit(0.86, 20, C.CString(modelPath), C.CString(logFilename))
	if ret != 1 {
		return errors.New("【Init】hiarClusterInit error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

func Registe(path string, registePictureFile []string) (failedPictureFile []string, err error) {

	l := len(registePictureFile)
	if l == 0 {
		l = 1
	}

	var failedNum C.int
	var failInfo = make([]C.ImageInfo, l)

	ret := C.hiarAddingImages(C.CString(path), &failInfo[0], &failedNum)
	if ret != 1 {
		return nil, errors.New("【Registe】hiarAddingImages error, retCode:" + strconv.Itoa(int(ret)))
	}

	fmt.Println("hiarAddingImages , failedNum:", failedNum)

	for i := 0; i < int(failedNum); i++ {
		fmt.Println("hiarAddingImages failed image: ", C.GoString(&failInfo[i].filename[0]))
		failedPictureFile = append(failedPictureFile, C.GoString(&failInfo[i].filename[0]))
	}

	return failedPictureFile, nil
}

func RegisteSingle(image *Image, filename string) (err error) {
	var cImage C.ImageData
	cImage.data = (*C.uchar)((unsafe.Pointer)(&image.Data[0]))
	cImage.data_len = C.int(len(image.Data))
	cImage.width = C.int(image.Width)
	cImage.height = C.int(image.Height)
	cImage.data_type = C.enum_ImageDataType(image.DataType)

	fmt.Println("filepath.Base(filename)", filepath.Base(filename))
	if ret := C.hiarAddingImage(&cImage, C.CString(filepath.Base(filename))); ret != OK {
		return errors.New("注册失败" + strconv.Itoa(int(ret)))
	}
	return nil
}

func Search2(image *Image) [FACE_MAX_RESULT]C.ImageInfo {

	//cImage := NewC_ImageData(image)

	var imageData C.ImageData
	//todo free

	imageData.data = (*C.uchar)((unsafe.Pointer)(&image.Data[0]))
	imageData.data_len = C.int(len(image.Data))
	imageData.width = C.int(image.Width)
	imageData.height = C.int(image.Height)
	imageData.data_type = C.enum_ImageDataType(image.DataType)

	var info = [FACE_MAX_RESULT]C.ImageInfo{}
	var v_len = C.int(FACE_MAX_RESULT)
	resultNum := C.hiarQuery(&imageData, &info[0], v_len)
	fmt.Println("hiarQuery resultNum ", resultNum, info[0])
	if int(resultNum) == 0 {
		return info
	}
	return info
}

func Search(image *Image) (results []*FaceEntity) {

	var imageData C.ImageData
	imageData.data = (*C.uchar)((unsafe.Pointer)(&image.Data[0]))
	imageData.data_len = C.int(len(image.Data))
	imageData.width = C.int(image.Width)
	imageData.height = C.int(image.Height)
	imageData.data_type = C.enum_ImageDataType(image.DataType)

	var info = make([]C.ImageInfo, FACE_MAX_RESULT)
	var v_len = C.int(FACE_MAX_RESULT)
	resultNum := C.hiarQuery(&imageData, &info[0], v_len)
	fmt.Println("hiarQuery resultNum ", resultNum, info[0])
	if int(resultNum) == 0 {
		return
	}

	//C.printImageInfo(&info[0], resultNum)

	for i, imageInfo := range info {
		if i < int(resultNum) {
			filename := C.GoString(&imageInfo.filename[0])
			match := float32(imageInfo.similarity)
			fmt.Println("【go range】 filename:", filename, "match:", match)

			results = append(results, &FaceEntity{
				RegFilename: C.GoString(&imageInfo.filename[0]),
				Match:       float32(imageInfo.similarity),
			})
		}

	}

	return
}

func UnRegisteAll() error {
	var imageInfo C.ImageInfo

	ret := C.hiarDelImages(&imageInfo, 0)
	if ret != 1 {
		return errors.New("【UnRegisteAll】hiarDelImages error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

func UnRegiste() {

	var cFilename [C.MAX_FILE_NAME_LEN]C.char

	cstr := C.CString("DSC08066.jpg")
	defer C.free(unsafe.Pointer(cstr))
	C.strncpy(&cFilename[0], cstr, C.MAX_FILE_NAME_LEN)

	var imageInfo C.ImageInfo
	imageInfo.filename = cFilename

	fmt.Println("注销人脸", C.GoString(&imageInfo.filename[0]))

	ret := C.hiarDelImages(&imageInfo, 1)
	fmt.Println("hiarDelImages ret", ret)
}
