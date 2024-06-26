package face

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/bighuangbee/face_search2/api/biz/v1"
	"github.com/bighuangbee/face_search2/app/internal/data"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var ErrorFaceRegistering = errors.New(400, "ErrorFaceRegistering", "Face registering, please wait")
var ErrorRequestMissingFile = errors.New(400, "ErrorRequestMissingFile", "请求缺少form file")
var ErrorImageTypeRequired = errors.New(400, "ErrorImageTypeRequired", "需要图片类型， "+strings.Join(face_wrapper.PictureExt, "、"))
var ErrorRequestFrom = errors.New(400, "ErrorRequestFrom", "ErrorRequestFrom")
var ErrorFaceSDK = errors.New(400, "ErrorFaceSDK", "ErrorFaceSDK")
var ErrorFaceSearchEmpty = errors.New(400, "ErrorFaceSearchEmpty", "ErrorFaceSearchEmpty")

type FaceRecognizeApp struct {
	logger      log.Logger
	log         *log.Helper
	data        *data.Data
	registering atomic.Bool
	bc          *conf.Bootstrap

	regService *RegisteService
}

const FACE_REGISTE_PATH = "/hiar_face/registe_path"
const FACE_REGISTE_LOGS = "/hiar_face/registe_logs"
const FACE_SEARCH_RECORD = "/hiar_face/search_record"
const FACE_REGISTE_FAILED = "/hiar_face/registe_failed"

var registeLogFile = FACE_REGISTE_LOGS + "/" + time.Now().Format("2006-01-02") + ".txt"
var searchRecordFilename = FACE_SEARCH_RECORD + "/" + time.Now().Format("2006-01-02") + ".txt"

func NewFaceRecognizeApp(logger log.Logger, bc *conf.Bootstrap, data *data.Data) *FaceRecognizeApp {
	app := FaceRecognizeApp{
		logger:     logger,
		log:        log.NewHelper(log.With(logger, "module", "service/FaceRecognizeApp")),
		data:       data,
		bc:         bc,
		regService: NewRegisteService(logger),
	}

	face_registe_path := os.Getenv("face_registe_path")
	if face_registe_path == "" {
		face_registe_path = FACE_REGISTE_PATH
	}

	face_models_path := os.Getenv("face_models_path")
	if face_models_path == "" {
		face_models_path = "/root/face_search/libs/models/"
	}

	os.MkdirAll(face_registe_path, 0755)
	os.MkdirAll(FACE_REGISTE_LOGS, 0755)
	os.MkdirAll(FACE_SEARCH_RECORD, 0755)
	os.MkdirAll(FACE_REGISTE_FAILED, 0755)

	app.log.Infow("face_registe_path", face_registe_path, "face_models_path", face_models_path)

	svcPath := bc.Face.SearchSvcPath
	if bc.Face.FaceMode == conf.FaceMode_registe {
		svcPath = bc.Face.RegisteSvcPath
	}
	os.MkdirAll(bc.Face.RegisteSvcPath, 0755)
	os.MkdirAll(bc.Face.SearchSvcPath, 0755)

	if err := face_wrapper.Init(face_models_path, bc.Face.GetMatch(), svcPath); err != nil {
		app.log.Infow("【NewFaceRecognizeApp】face_wrapper init", err)
		panic(err)
	}

	return &app
}

func (s *FaceRecognizeApp) RegisteByPath(ctx context.Context, req *pb.RegisteRequest) (*pb.RegisteByPathReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	registedSuccFace, registedFailedFace, newFace, _ := RegFilePreProcess()
	s.log.Infow("本次新增人脸", len(newFace), "之前注册成功人脸", len(registedSuccFace), "之前注册失败的人脸", len(registedFailedFace))

	fn := func() {
		s.registering.Store(true)
		defer s.registering.Store(false)
		//s.registeFace()
		s.registeFaceOneByOne(registedSuccFace, newFace, false)
	}

	if req.GetSync() {
		//阻塞运行耗时久，不适合http调用
		fn()
	} else {
		go fn()
	}

	return &pb.RegisteByPathReply{
		RegistedSuccNum:   int32(len(registedSuccFace)),
		RegistedFailedNum: int32(len(registedFailedFace)),
		NewFaceNum:        int32(len(newFace)),
	}, nil
}

func (s *FaceRecognizeApp) RegisteStatus(context.Context, *pb.EmptyRequest) (*pb.RegisteStatusReply, error) {
	return &pb.RegisteStatusReply{
		Registering: s.registering.Load(),
	}, nil
}

func (s *FaceRecognizeApp) UnRegisteAll(ctx context.Context, req *pb.EmptyRequest) (*pb.NotifyReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	os.Remove(registeLogFile)

	s.regService.Repo.Range(func(key, value interface{}) bool {
		s.regService.Repo.Delete(key)
		return true
	})

	err := face_wrapper.UnRegisteAll()
	if err != nil {
		return nil, ErrorFaceSDK
	}

	resp, err := util.HttpPost("http://localhost:6666/unregiste", map[string]interface{}{})

	s.log.Log(log.LevelInfo, "通知【注册服务】", "", err, string(resp))

	return &pb.NotifyReply{Ok: true}, nil
}

type SearchRecord struct {
	Time      string                     `json:"time"`
	Filename  string                     `json:"filename"`
	StartTime string                     `json:"startTime"`
	EndTime   string                     `json:"endTime"`
	Results   []*face_wrapper.FaceEntity `json:"results"`
}

func (s *FaceRecognizeApp) Search(ctx context.Context) (reply *pb.SearchResultReply, err error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	request, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return nil, ErrorRequestFrom
	}

	image, filename, err := receiveFaceFile(request)
	if err != nil {
		return nil, err
	}

	results := face_wrapper.Search(image)
	reply = &pb.SearchResultReply{}

	startTime := time.Time{}
	endTime := time.Time{}
	if len(results) > 0 && s.bc.Face.GetMatchTimeRange() > 0 {
		t, _ := GetShootTime(results[0].RegFilename)
		startTime = t.Add(time.Duration(-s.bc.Face.GetMatchTimeRange()) * time.Minute)
		endTime = t.Add(time.Duration(s.bc.Face.GetMatchTimeRange()) * time.Minute)

		fmt.Println("GetShootTime", t.String(), "startTime", startTime.String(), "endTime", endTime.String())
	}
	for _, result := range results {
		t, _ := GetShootTime(result.RegFilename)
		if s.bc.Face.GetMatchTimeRange() > 0 {
			if !(t.After(startTime) && t.Before(endTime)) {
				fmt.Printf(" 排除result.RegFilename", result.RegFilename)
				continue
			}
		}

		reply.Results = append(reply.Results, &pb.SearchResult{
			Filename:  result.RegFilename,
			Match:     result.Match,
			ShootTime: t.Format("01-02 15:04:05"),
		})

	}

	//startTime := request.FormValue("startTime")
	//endTime := request.FormValue("endTime")

	//s.log.Infow("Search formdata", "", "inputTimeStr", startTime, "endTime", endTime)

	////算法搜索不到结果时，按时间范围检索图片
	//if len(results) == 0 && startTime != "" && endTime != "" {
	//	fileInforesults, err := GetRangeFile(s.FileInfoRepo, startTime, endTime)
	//	if err != nil {
	//		s.log.Errorw("GetRangeFile", err)
	//		return nil, ErrorRequestFrom
	//	}
	//	for _, result := range fileInforesults {
	//		reply.Results = append(reply.Results, &pb.SearchResult{
	//			Filename: result.Filename,
	//		})
	//	}
	//
	//	s.log.Infow("算法检索不到结果, 进行文件时间检索, 结果数量:", len(fileInforesults), "fileInforesults", fileInforesults)
	//}

	if len(reply.Results) == 0 {
		err = ErrorFaceSearchEmpty
	}

	go func() {
		basePath := FACE_SEARCH_RECORD + "/" + time.Now().Format("2006-01-02") + "/"
		os.MkdirAll(basePath, 0755)
		os.WriteFile(basePath+filename, image.Data, 0644)

		str, _ := json.Marshal(&SearchRecord{
			Time:     util.GetLocTime(),
			Filename: basePath + filename,
			Results:  results,
		})
		if err := util.CreateOrOpenFile(searchRecordFilename, string(str)); err != nil {
			s.log.Errorw("Search CreateOrOpenFile", err)
		}

	}()

	return reply, err
}

func (s *FaceRecognizeApp) FaceSearchByDatetime(ctx context.Context, req *pb.FaceSearchByDatetimeRequest) (reply *pb.SearchResultReply, err error) {

	searchRecord := SearchRecord{
		Time:      util.GetLocTime(),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	reply = &pb.SearchResultReply{}

	//按时间范围检索照片
	if req.StartTime != "" && req.EndTime != "" {
		fileInforesults, err := GetRangeFile(s.regService.Repo, req.StartTime, req.EndTime)
		if err != nil {
			s.log.Errorw("GetRangeFile", err)
			return nil, ErrorRequestFrom
		}
		for _, result := range fileInforesults {
			reply.Results = append(reply.Results, &pb.SearchResult{
				Filename:  result.Filename,
				ShootTime: result.ShootTime.Format("01-02 15:04:05"),
			})

			searchRecord.Results = append(searchRecord.Results, &face_wrapper.FaceEntity{
				RegFilename: result.Filename,
			})
		}

		s.log.Infow("按照片的生成时间搜索, 结果数量:", len(fileInforesults), "fileInforesults", fileInforesults)
	}

	if len(reply.Results) == 0 {
		err = ErrorFaceSearchEmpty
	}

	go func() {
		basePath := FACE_SEARCH_RECORD + "/" + time.Now().Format("2006-01-02") + "/"
		os.MkdirAll(basePath, 0755)

		str, _ := json.Marshal(&searchRecord)
		if err := util.CreateOrOpenFile(searchRecordFilename, string(str)); err != nil {
			s.log.Errorw("FaceSearchByDatetime CreateOrOpenFile", err)
		}
	}()

	return reply, err
}

func (s *FaceRecognizeApp) FaceDbReload(ctx context.Context, req *pb.EmptyRequest) (reply *pb.NotifyReply, err error) {
	if err = copyFile(s.bc.Face.RegisteSvcPath+"/"+face_wrapper.DbName, s.bc.Face.SearchSvcPath+"/"+face_wrapper.DbName); err != nil {
		s.log.Infow("复制db文件失败", err)
	}

	if err = face_wrapper.LoadDB(s.bc.Face.SearchSvcPath); err != nil {
		s.log.Infow("SDK加载db失败", err)
	}

	fileList, err1 := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err1 != nil {
		s.log.Infow("GetFilesWithExtensions", err)
		err = err1
	} else {
		s.log.Info("加载RegisteInfo")
		for _, filename := range fileList {
			t, _ := GetShootTime(filename)
			s.regService.Repo.Store(filename, face_wrapper.RegisteInfo{
				Filename:  filename,
				ShootTime: t,
			})
		}

	}

	s.log.Info("【搜索服务】db复制并加载成功")
	return &pb.NotifyReply{Ok: true}, nil
}
