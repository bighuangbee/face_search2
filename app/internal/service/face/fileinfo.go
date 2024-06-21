package face

import (
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sys/unix"
	"time"
)

type FileInfo struct {
	Filename  string
	Birthtime time.Time
}

// 加载时区
var location, _ = time.LoadLocation("Asia/Shanghai")
var timeFormat = "2006-01-02 15:04"

func LoadFileInfo() map[string]*FileInfo {
	//fileList, err := util.GetFilesWithExtensions("../../../../libs/data/gallery", []string{".png", ".jpg", ".jpeg"})
	fileList, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err != nil {
		log.Errorw("【RegisteByPath】GetFilesWithExtensions", err)
		return nil
	}

	var FileInfoRepo = make(map[string]*FileInfo, 0)
	for _, filename := range fileList {

		var stat unix.Stat_t
		if err := unix.Stat(filename, &stat); err != nil {
			fmt.Printf("获取文件系统信息时出错: %v\n", err)
			continue
		}

		//stat := fileInfo.Sys().(*syscall.Stat_t)

		fmt.Println("LoadFileInfo", filename, "Ctim:", time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).String(), "Mtim:", time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).String(), "Atim:", time.Unix(stat.Atim.Sec, stat.Atim.Nsec).String())
		//fmt.Println("Birthtimespec", filename, time.Unix(stat.Btim.Sec, stat.Btim.Nsec).String())

		FileInfoRepo[filename] = &FileInfo{
			Filename:  filename,
			Birthtime: time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).In(location),
		}
	}
	return FileInfoRepo
}

func GetRangeFile(fileInfoList map[string]*FileInfo, inputTimeStr string, timeRange time.Duration) (results []*FileInfo, err error) {
	inputTime, err := time.ParseInLocation(timeFormat, inputTimeStr, location)
	if err != nil {
		return nil, err
	}

	// 计算时间范围的开始和结束时间
	startTime := inputTime.Add(-timeRange)
	endTime := inputTime.Add(timeRange)

	// 查找在指定时间范围内的文件
	for _, info := range fileInfoList {
		if info.Birthtime.After(startTime) && info.Birthtime.Before(endTime) {
			results = append(results, info)
		}
	}

	return results, nil
}
