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
	"path/filepath"
	"strconv"
	"unsafe"
)

const FACE_MAX_RESULT = 20
const OK = C.int(1)

var PictureExt = []string{".png", ".jpg", ".jpeg"}

type RegisteInfo struct {
	Time     string `json:"time"`
	Ok       bool   `json:"ok"`
	Filename string `json:"filename"`
}

func Init(modelPath string, logFilename string, match float32) error {
	ret := C.hiarClusterInit(C.float(match), 20, C.CString(modelPath), C.CString(logFilename))
	if ret != 1 {
		return errors.New("【Init】hiarClusterInit error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

func Registe(regInfo []*RegisteInfo) (err error) {
	if len(regInfo) == 0 {
		return errors.New("Input regInfo empty.")
	}

	var inputList = make([]C.ImageInfo, len(regInfo))
	var succList = make([]C.ImageInfo, len(regInfo))

	for _, info := range regInfo {
		inputList = append(inputList, C.ImageInfo{
			filename: toCString(info.Filename),
		})
	}

	okNum := C.hiarAddingImages(&inputList[0], C.int(len(inputList)), &succList[0], C.int(len(regInfo)))
	if okNum < 0 {
		return errors.New("【Registe】hiarAddingImages error, retCode:" + strconv.Itoa(int(okNum)))
	}
	return nil
}

func RegisteSingle(image *Image, filename string) (err error) {
	if len(image.Data) == 0 {
		return errors.New("空照片")
	}
	var cImage C.ImageData
	cImage.data = (*C.uchar)((unsafe.Pointer)(&image.Data[0]))
	cImage.data_len = C.int(len(image.Data))
	cImage.width = C.int(image.Width)
	cImage.height = C.int(image.Height)
	cImage.data_type = C.enum_ImageDataType(image.DataType)

	if ret := C.hiarAddingImage(&cImage, C.CString(filename)); ret != OK {
		return errors.New("注册失败" + strconv.Itoa(int(ret)))
	}
	return nil
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
	if int(resultNum) == 0 {
		return
	}

	//C.printImageInfo(&info[0], resultNum)
	for i, imageInfo := range info {
		if i < int(resultNum) {
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
	if ret := C.hiarDelImages(&imageInfo, 0); ret != 1 {
		return errors.New("【UnRegisteAll】hiarDelImages error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

func UnRegiste(filename string) error {
	ret := C.hiarDelImages(&C.ImageInfo{
		filename: toCString(filepath.Base(filename)),
	}, 1)
	if ret != 1 {
		return errors.New("【UnRegiste】hiarDelImages error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

// toCString 将Go字符串转换为[C.MAX_FILE_NAME_LEN]C.char数组
func toCString(str string) [C.MAX_FILE_NAME_LEN]C.char {
	var cFilename [C.MAX_FILE_NAME_LEN]C.char
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.strncpy(&cFilename[0], cstr, C.MAX_FILE_NAME_LEN)
	return cFilename
}
