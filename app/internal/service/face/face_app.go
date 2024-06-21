package face

import (
	"context"
	"encoding/json"
	pb "github.com/bighuangbee/face_search2/api/biz/v1"
	"github.com/bighuangbee/face_search2/app/internal/data"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io/ioutil"
	"os"
	"strconv"
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
	log          *log.Helper
	data         *data.Data
	registering  atomic.Bool
	FileInfoRepo map[string]*FileInfo
}

const FACE_REGISTE_PATH = "/hiar_face/registe_path"
const FACE_REGISTE_LOGS = "/hiar_face/registe_logs"
const FACE_SEARCH_RECORD = "/hiar_face/search_record"

var logFilename = FACE_REGISTE_LOGS + "/" + time.Now().Format("2006-01-02") + ".txt"
var searchRecordFilename = FACE_SEARCH_RECORD + "/" + time.Now().Format("2006-01-02") + ".txt"

func NewFaceRecognizeApp(logger log.Logger, bc *conf.Bootstrap, data *data.Data) *FaceRecognizeApp {

	app := FaceRecognizeApp{
		log:  log.NewHelper(log.With(logger, "module", "service/FaceRecognizeApp")),
		data: data,
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

	app.log.Infow("face_registe_path", face_registe_path, "face_models_path", face_models_path)

	if err := face_wrapper.Init(face_models_path, "./hiarClusterLog.txt", bc.Match); err != nil {
		app.log.Infow("【NewFaceRecognizeApp】face_wrapper init", err)
		panic(err)
	}

	if bc.StartUnRegiste {
		if err := face_wrapper.UnRegisteAll(); err != nil {
			app.log.Warnw("【NewFaceRecognizeApp】UnRegisteAll ", err)
		}
	}

	registedFace, _, newFace := facePreProcessing(app.log)
	app.registeFaceOneByOne(registedFace, newFace, true)

	if bc.RegisteTimer > 0 {
		ticker := time.NewTicker(time.Minute * time.Duration(bc.RegisteTimer))
		go func() {
			defer ticker.Stop()
			for {
				<-ticker.C
				if !app.registering.Load() {
					app.registering.Store(true)
					app.log.Infow("定时自动注册人脸", "")
					registedFace, _, newFace := facePreProcessing(app.log)
					app.registeFaceOneByOne(registedFace, newFace, true)
					app.registering.Store(false)
				} else {
					app.log.Infow("定时自动注册人脸失败", "上次注册执行中")
				}

			}
		}()
	}

	app.FileInfoRepo = LoadFileInfo()
	return &app
}

func (s *FaceRecognizeApp) RegisteByPath(context.Context, *pb.EmptyRequest) (*pb.RegisteByPathReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	registedSuccFace, registedFailedFace, newFace := facePreProcessing(s.log)
	s.log.Infow("本次新增人脸", len(newFace), "之前注册成功人脸", len(registedSuccFace), "之前注册失败的人脸", len(registedFailedFace))

	go func() {
		s.registering.Store(true)
		defer s.registering.Store(false)
		//s.registeFace()
		s.registeFaceOneByOne(registedSuccFace, newFace, false)
	}()

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

func (s *FaceRecognizeApp) UnRegisteAll(ctx context.Context, req *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	os.Remove(logFilename)

	err := face_wrapper.UnRegisteAll()
	if err != nil {
		return nil, ErrorFaceSDK
	}

	return &pb.EmptyReply{}, nil
}

type SearchRecord struct {
	Time     string                     `json:"time"`
	Filename string                     `json:"filename"`
	Results  []*face_wrapper.FaceEntity `json:"results"`
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
	for _, result := range results {
		reply.Results = append(reply.Results, &pb.SearchResult{
			Filename: result.RegFilename,
			Match:    result.Match,
		})
	}

	inputTimeStr := request.FormValue("inputTime")
	timeRangeValue, _ := strconv.Atoi(request.FormValue("timeRange"))
	timeRange := time.Minute * time.Duration(timeRangeValue)

	s.log.Infow("Search form data", inputTimeStr, timeRangeValue)

	//算法搜索不到结果时，按时间范围检索图片
	if len(results) == 0 {

		fileInforesults, err := GetRangeFile(s.FileInfoRepo, inputTimeStr, timeRange)
		if err != nil {
			s.log.Errorw("GetRangeFile", err)
			return nil, ErrorRequestFrom
		}
		for _, result := range fileInforesults {
			reply.Results = append(reply.Results, &pb.SearchResult{
				Filename: result.Filename,
			})
		}
	}

	if len(reply.Results) == 0 {
		err = ErrorFaceSearchEmpty
	}

	go func() {
		basePath := FACE_SEARCH_RECORD + "/" + time.Now().Format("2006-01-02") + "/"
		os.MkdirAll(basePath, 0755)

		ioutil.WriteFile(basePath+filename, image.Data, 0644)

		str, _ := json.Marshal(&SearchRecord{
			Time:     util.GetLocTime(),
			Filename: basePath + filename,
			Results:  results,
		})
		if err := util.CreateOrOpenFile(searchRecordFilename, string(str)); err != nil {
			s.log.Errorw("CreateOrOpenFile", err)
		}

	}()

	return reply, err
}
