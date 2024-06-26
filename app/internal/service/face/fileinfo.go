package face

import (
	"errors"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"sync"
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
	//////unix获取不到Btim
	//var stat unix.Stat_t
	//if err := unix.Stat(filename, &stat); err != nil {
	//	return time.Time{}, err
	//}
	//
	//return time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).In(location), nil

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return time.Time{}, err
	}
	return x.DateTime()

}

// GetRangeFile 查找在指定时间范围内的文件
func GetRangeFile(fileInfoList sync.Map, startTimeStr string, endTimeStr string) (results []*face_wrapper.RegisteInfo, err error) {
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
	fileInfoList.Range(func(key, value any) bool {
		info := value.(face_wrapper.RegisteInfo)
		if info.ShootTime.After(startTime) && info.ShootTime.Before(endTime) {
			results = append(results, &info)
		}
		return true
	})

	return results, nil
}
