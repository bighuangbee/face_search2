package face

import (
	"encoding/json"
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/app/internal/service/storage"
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RegisteService struct {
	logger log.Logger
	FaceDb storage.FaceDb
	queue  chan *face_wrapper.RegisteInfo
	bc     *conf.Bootstrap
}

func NewRegisteService(logger log.Logger, config *conf.Bootstrap) *RegisteService {
	s := &RegisteService{
		logger: logger,
		queue:  make(chan *face_wrapper.RegisteInfo, 2000),
		bc:     config,
	}

	fmt.Println("NewRegisteService 1")
	var err error
	if s.FaceDb, err = storage.NewMysqlStorage(config.Data, logger); err != nil {
		logger.Log(log.LevelError, "storage.NewBoltDb", err)
		panic(err)
	}
	fmt.Println("NewRegisteService 2")
	s.logger.Log(log.LevelInfo, "face bb open success.")

	return s
}

func (s *RegisteService) PushQueue(filename string) {
	s.queue <- &face_wrapper.RegisteInfo{
		Filename: filename,
	}
}

func (s *RegisteService) Run() {
	for {
		info := <-s.queue
		_, ok := s.FaceDb.Read(info.Filename)
		if ok {
			continue
		}

		t1 := time.Now()
		regInfo, err := s.Reg(info.Filename)
		if err != nil {
			s.logger.Log(log.LevelError, "注册出错", info.Filename, "耗时", time.Since(t1))
			continue
		}

		s.logger.Log(log.LevelInfo, "注册结果 ok", regInfo.Ok, "fielname", info.Filename, "耗时", time.Since(t1))
	}
}

func (s *RegisteService) UnReg(filename string) {
	if err := face_wrapper.UnRegiste(filename); err != nil {
		s.logger.Log(log.LevelError, "face_wrapper.UnRegiste error", err, "filename", filename)
	}
	if err := s.FaceDb.Delete(filename); err != nil {
		s.logger.Log(log.LevelError, "regService.FaceDb.Delete error", err, "filename", filename)
	}
	if err := os.Remove(filename); err != nil {
		s.logger.Log(log.LevelError, "os.Remove error", err, "filename", filename)
	}
}

func (s *RegisteService) Reg(filename string) (regInfo *storage.RegisteInfo, err error) {
	imageFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regError := face_wrapper.RegisteSingle(&face_wrapper.Image{
		DataType: face_wrapper.GetImageType(filename),
		Size:     len(imageFile),
		Data:     imageFile,
	}, filename)

	bTime, _ := GetShootTime(filename)
	if bTime.IsZero() {
		bTime, _ = GetCreateTime(filename)
	}

	regInfo = &storage.RegisteInfo{
		Filename:  filename,
		Ok:        false,
		Time:      util.GetLocTime(),
		ShootTime: bTime,
	}

	if regError == nil {
		regInfo.Ok = true
		if err := s.FaceDb.Update(regInfo.Filename, regInfo); err != nil {
			s.logger.Log(log.LevelError, "FaceDb.Update", err)
		}

	} else {
		//注册失败的照片
		os.Rename(filename, FACE_REGISTE_FAILED+"/"+time.Now().Format("20060102")+"_"+filepath.Base(filename))
	}

	str, _ := json.Marshal(regInfo)
	util.CreateOrOpenFile(registeDataLogsDay, string(str))

	return regInfo, nil
}

func (s RegisteService) CheckExpired() (values []*storage.RegisteInfo) {
	//检查图片文件是否过期
	fileList, _ := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	for _, filename := range fileList {
		if strings.HasPrefix(filepath.Base(filename), "test_") {
			return
		}
		st, _ := GetShootTime(filename)
		isExpired := s.IsExpired(st)
		if isExpired {
			s.UnReg(filename)
		}

		s.logger.Log(log.LevelInfo, "CheckExpired, photo isExpired:", isExpired, "filename", filename, "ShootTime", st.String())
	}

	//检查特征库是否过期
	expiredList, err := s.FaceDb.DeleteExpired(time.Duration(s.bc.Face.EffectiveTime) * time.Hour)
	if err != nil {
		s.logger.Log(log.LevelError, "FaceDb.DeleteExpired", err)
		return
	}
	s.logger.Log(log.LevelInfo, "FaceDb, 清理特征库过期数量", len(expiredList))
	return expiredList
}

func (s RegisteService) IsExpired(t time.Time) bool {
	//todo
	effectiveDuration := time.Duration(s.bc.Face.EffectiveTime) * time.Hour
	cutoffTime := time.Now().Add(-effectiveDuration).In(location)
	//fmt.Println("p------", cutoffTime.Format(time.DateTime), t.Format(time.DateTime), t.Before(cutoffTime), cutoffTime.Before(t))

	return t.Before(cutoffTime)

}
