package face

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

/**
 * @Desc  人脸注册图预处理
 * @return 已注册的人脸，待注册的新人脸
 **/
func RegFilePreProcess() (registedSuccFace []string, registedFailedFace []string, newFace []string, err error) {

	//对容器内注册文件重命名
	//files, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	//if err != nil {
	//	log.Errorw("【RegisteByPath】GetFilesWithExtensions", err)
	//	return
	//}
	//
	//renameFlag := "_"
	//fileflag := "_" + time.Now().Format("01021504") + renameFlag
	//for _, filename := range files {
	//	if !strings.HasSuffix(filename, renameFlag+filepath.Ext(filename)) {
	//		rename := filepath.Dir(filename) + "/" + util.GetFileName(filename) + fileflag + filepath.Ext(filename)
	//		if err := os.Rename(filename, rename); err != nil {
	//			log.Warnf("os.Rename", err)
	//		}
	//		//log.Infow("注册图重命名", "", "filename", filename, "rename", rename)
	//	}
	//}

	fileList, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err != nil {
		return nil, nil, nil, err
	}

	//之前已注册的人脸
	registedFaceMap := make(map[string]*face_wrapper.RegisteInfo)

	file, err := os.Open(registeLogFile)
	if err != nil {
		return []string{}, []string{}, fileList, err
	}
	defer file.Close()

	//检查注册日志，不重复注册
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := face_wrapper.RegisteInfo{}
		json.Unmarshal(scanner.Bytes(), &line)
		registedFaceMap[line.Filename] = &line
	}

	if err := scanner.Err(); err != nil {
		return []string{}, []string{}, fileList, err
	}

	for _, filename := range fileList {
		if item, ok := registedFaceMap[filename]; ok {
			if item.Ok {
				registedSuccFace = append(registedSuccFace, filename)
			} else {
				registedFailedFace = append(registedFailedFace, filename)
			}
		} else {
			newFace = append(newFace, filename)
		}
	}
	return
}

/**
 * @Desc  逐个图片注册
 * @Param reset 是否跳过已注册的图片
 **/
func (s *FaceRecognizeApp) registeFaceOneByOne(registedFace []string, newFace []string, reset bool) {

	t := time.Now()
	s.log.Infow("【registeFaceOneByOne】人脸注册开始", "")

	if reset {
		newFace = append(newFace, registedFace...)
	}

	//注册的图片数量
	registeNum := len(newFace)
	//注册成功的数量
	registeSuccNum := 0

	for index, filename := range newFace {
		t1 := time.Now()
		result, err := s.regService.Reg(filename)
		if err == nil {
			s.log.Errorw("注册出错", strconv.Itoa(index+1)+" "+filename, "耗时", time.Since(t1))
			continue
		}
		if result.Ok {
			s.log.Infow("注册成功", strconv.Itoa(index+1)+" "+filename, "耗时", time.Since(t1))
		} else {
			s.log.Infow("注册失败", strconv.Itoa(index+1)+" "+filename, "耗时", time.Since(t1))
		}
	}

	s.log.Infow("【registeFaceOneByOne】人脸注册结束", "", "新增注册人脸数量", registeNum, "注册成功", registeSuccNum, "注册失败", registeNum-registeSuccNum, "耗时", time.Since(t))
}

func (s *FaceRecognizeApp) registeFace() {

	t := time.Now()
	s.log.Infow("【RegisteByPath】begining", "")

	files, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err != nil {
		s.log.Errorw("【RegisteByPath】GetFilesWithExtensions", err)
		return
	}

	var regInfo []*face_wrapper.RegisteInfo
	for _, filename := range files {
		regInfo = append(regInfo, &face_wrapper.RegisteInfo{
			Filename: filename,
		})
	}

	err = face_wrapper.Registe(regInfo)
	if err != nil {
		s.log.Errorw("【RegisteByPath】failed", err)
	} else {
		s.log.Infow("【RegisteByPath】end", "success", "duration", time.Since(t))
	}

	for index, info := range regInfo {
		result := face_wrapper.RegisteInfo{
			Filename: info.Filename,
			Ok:       info.Ok,
			Time:     util.GetLocTime(),
		}

		str, _ := json.Marshal(&result)
		if err := util.CreateOrOpenFile(registeLogFile, string(str)); err != nil {
			s.log.Errorw("CreateOrOpenFile", err)
		}

		if result.Ok {
			s.log.Infow("注册成功", strconv.Itoa(index+1)+" "+info.Filename)
		} else {
			s.log.Infow("注册失败", strconv.Itoa(index+1)+" "+info.Filename)
		}
	}

}

func receiveFaceFile(request *http.Request) (image *face_wrapper.Image, filename string, err error) {
	file, fileHeader, err := request.FormFile("file")

	if err != nil {
		fmt.Println("err", err)
		return nil, "", ErrorRequestMissingFile
	}
	defer file.Close()

	if !util.HasValidExtension(filepath.Ext(fileHeader.Filename), face_wrapper.PictureExt) {
		return nil, "", ErrorImageTypeRequired
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	image = &face_wrapper.Image{
		Data:     fileData,
		Size:     len(fileData),
		DataType: face_wrapper.GetImageType(fileHeader.Filename),
	}

	filename = fileHeader.Filename
	return
}

// copyFile 复制文件，如果目标文件已存在则覆盖
func copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcF.Close()

	// 打开（创建或覆盖）目标文件
	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstF.Close()

	// 将源文件内容复制到目标文件
	_, err = io.Copy(dstF, srcF)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// 刷新写入的内容
	err = dstF.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
