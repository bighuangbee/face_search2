package face

import (
	"errors"
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/app/internal/service/storage"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/sys/unix"
	"os"
	"time"
)

type FileInfo struct {
	Filename  string
	Birthtime time.Time
}

// 加载时区
var location, _ = time.LoadLocation("Asia/Shanghai")
var timeFormat = "2006-01-02 15:04"

func LoadFileInfo(logger log.Logger) map[string]*FileInfo {
	//fileList, err := util.GetFilesWithExtensions("../../../../libs/data/gallery", []string{".png", ".jpg", ".jpeg"})
	fileList, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err != nil {
		log.Errorw("【RegisteByPath】GetFilesWithExtensions", err)
		return nil
	}

	var FileInfoRepo = make(map[string]*FileInfo, 0)
	for _, filename := range fileList {
		t, _ := GetShootTime(filename)
		FileInfoRepo[filename] = &FileInfo{
			Filename:  filename,
			Birthtime: t,
		}

		logger.Log(log.LevelInfo, "filename", filename, "GetShootTime", t.Format(time.DateTime))
	}
	return FileInfoRepo
}

func GetShootTime(filename string) (time.Time, error) {
	file, err := os.Open(filename)
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return time.Time{}, err
	}
	return x.DateTime()
}

func GetCreateTime(filename string) (time.Time, error) {
	var stat unix.Stat_t
	if err := unix.Stat(filename, &stat); err != nil {
		return time.Time{}, err
	}
	return time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).In(location), nil
}

// GetRangeFile 查找在指定时间范围内的文件
func GetRangeFile(fileInfoList []*storage.RegisteInfo, startTimeStr string, endTimeStr string) (results []*face_wrapper.RegisteInfo, err error) {
	fmt.Println("GetRangeFile ", startTimeStr, endTimeStr)
	startTime, err := time.Parse(timeFormat, startTimeStr)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(timeFormat, endTimeStr)
	if err != nil {
		return nil, err
	}

	// 确保开始时间在结束时间之前
	if startTime.After(endTime) {
		return nil, errors.New("startTime should be before endTime")
	}

	// 查找在指定时间范围内的文件
	for _, value := range fileInfoList {
		info := (*face_wrapper.RegisteInfo)(value)
		fmt.Println("fileInfoList range  ", value.Filename, value.ShootTime.String(), info.ShootTime.After(startTime) && info.ShootTime.Before(endTime))
		if info.ShootTime.After(startTime) && info.ShootTime.Before(endTime) {
			results = append(results, info)
		}
	}

	return results, nil
}
