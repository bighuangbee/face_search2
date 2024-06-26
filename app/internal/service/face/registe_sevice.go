package face

import (
	"encoding/json"
	"errors"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RegisteService struct {
	logger log.Logger
	Repo   sync.Map
	queue  chan *face_wrapper.RegisteInfo
}

func NewRegisteService(logger log.Logger) *RegisteService {
	s := &RegisteService{
		logger: logger,
		queue:  make(chan *face_wrapper.RegisteInfo, 2000),
	}

	fileData, err := os.ReadFile(registeLogFile)
	if err != nil {
		logger.Log(log.LevelInfo, "open registeLogFile", registeLogFile, "err", err)
	}

	data := []face_wrapper.RegisteInfo{}
	if err := json.Unmarshal(fileData, &data); err != nil {
		logger.Log(log.LevelInfo, "json.Unmarshal", err)
	}

	for _, item := range data {
		s.Repo.Store(item.Filename, item)
	}

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
		_, ok := s.Repo.Load(info.Filename)
		if ok {
			continue
		}

		t1 := time.Now()
		regInfo, err := s.Reg(info.Filename)
		if err != nil {
			s.logger.Log(log.LevelError, "注册出错", info.Filename, "耗时", time.Since(t1))
			continue
		}

		s.Repo.Store(info.Filename, info)

		s.logger.Log(log.LevelInfo, "注册结果 ok", regInfo.Ok, "fielname", info.Filename, "耗时", time.Since(t1))
	}
}

func (s *RegisteService) Reg(filename string) (regInfo *face_wrapper.RegisteInfo, err error) {
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
	regInfo = &face_wrapper.RegisteInfo{
		Filename:  filename,
		Ok:        false,
		Time:      util.GetLocTime(),
		ShootTime: bTime,
	}

	if regError == nil {
		regInfo.Ok = true
	} else {
		os.Rename(filename, FACE_REGISTE_FAILED+"/"+time.Now().Format("20060102")+"_"+filepath.Base(filename))
	}

	str, _ := json.Marshal(&regInfo)
	if err := util.CreateOrOpenFile(registeLogFile, string(str)); err != nil {
		return nil, errors.New("CreateOrOpenFile " + err.Error())
	}
	return regInfo, nil
}
