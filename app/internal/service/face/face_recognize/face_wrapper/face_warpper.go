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
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/bighuangbee/face_search2/pkg/util"
	"strconv"
	"unsafe"
)

const FACE_MAX_RESULT = 20

var PictureExt = []string{".png", ".jpg", ".jpeg"}

func Init(modelPath string, logFilename string) error {
	ret := C.hiarClusterInit(0.86, 20, C.CString(modelPath), C.CString(logFilename))
	if ret != 1 {
		return errors.New("【Init】hiarClusterInit error, retCode:" + strconv.Itoa(int(ret)))
	}
	return nil
}

func Registe(path string) (regNumber int, err error) {

	files, err := util.GetFilesWithExtensions(path, PictureExt)
	if err != nil {
		return 0, err
	}

	var failedNum C.int
	var failInfo = make([]C.ImageInfo, len(files))

	ret := C.hiarAddingImages(C.CString(path), &failInfo[0], &failedNum)
	if ret != 1 {
		return 0, errors.New("【Registe】hiarAddingImages error, retCode:" + strconv.Itoa(int(ret)))
	}

	fmt.Println("failInfo , failedNum:", failedNum)

	for i := 0; i < int(failedNum); i++ {
		fmt.Println("failInfo failed image: ", C.GoString(&failInfo[i].filename[0]))
	}

	return int(failedNum), nil
}

func Search(image *Image) (results []*FaceEntity) {
	cImage := NewC_ImageData(image)

	var info [FACE_MAX_RESULT]C.ImageInfo
	var v_len = C.int(FACE_MAX_RESULT)

	resultNum := C.hiarQuery(cImage, &info[0], v_len)
	if int(resultNum) == 0 {
		return
	}

	for i := 0; i < int(resultNum); i++ {
		fmt.Println(C.GoString(&info[i].filename[0]), float32(info[i].similarity))
		results = append(results, &FaceEntity{
			RegFilename: C.GoString(&info[i].filename[0]),
			Match:       float32(info[i].similarity),
		})
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
