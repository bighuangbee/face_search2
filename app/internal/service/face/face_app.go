package face

import (
	"bufio"
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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	log         *log.Helper
	data        *data.Data
	registering atomic.Bool
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

	if err := face_wrapper.Init(face_models_path, "./hiarClusterLog.txt"); err != nil {
		app.log.Infow("【NewFaceRecognizeApp】face_wrapper init", err)
		panic(err)
	}

	if err := face_wrapper.UnRegisteAll(); err != nil {
		app.log.Warnw("【NewFaceRecognizeApp】UnRegisteAll ", err)
	}

	app.registeFaceOneByOne(true)

	return &app
}

func (s *FaceRecognizeApp) RegisteByPath(context.Context, *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	go func() {
		s.registering.Store(true)
		defer s.registering.Store(false)
		//s.registeFace()
		s.registeFaceOneByOne(false)
	}()

	return &pb.EmptyReply{}, nil
}

/**
 * @Desc  逐个图片注册
 * @Param reset 是否跳过已注册的图片
 **/
func (s *FaceRecognizeApp) registeFaceOneByOne(reset bool) {

	t := time.Now()
	s.log.Infow("【registeFaceOneByOne】begining", "")

	files, err := util.GetFilesWithExtensions(FACE_REGISTE_PATH, face_wrapper.PictureExt)
	if err != nil {
		s.log.Errorw("【RegisteByPath】GetFilesWithExtensions", err)
		return
	}

	//之前已注册的人脸
	registedFaceMap := make(map[string]struct{})

	file, err := os.Open(logFilename)
	if err != nil {
		s.log.Warnw("os.Open error", err, "logFilename", logFilename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := face_wrapper.RegisteInfo{}
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			s.log.Errorw("json.Unmarshal", err)
		}
		registedFaceMap[line.Filename] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		s.log.Warnw("canner registe file error", err, "logFilename", logFilename)
	}

	//注册的图片数量
	registeNum := 0
	//注册成功的数量
	registeSuccNum := 0
	for index, filename := range files {
		if _, ok := registedFaceMap[filename]; ok && !reset {
			s.log.Infow("人脸已注册", strconv.Itoa(index+1)+" "+filename)
			continue
		}

		imageFile, err := ioutil.ReadFile(filename)
		if err != nil {
			s.log.Infow("ReadFile error", filename)
			continue
		}
		regError := face_wrapper.RegisteSingle(&face_wrapper.Image{
			DataType: face_wrapper.GetImageType(filename),
			Size:     len(imageFile),
			Data:     imageFile,
		}, filename)

		result := face_wrapper.RegisteInfo{
			Filename: filename,
			Ok:       false,
			Time:     util.GetLocTime(),
		}

		if regError == nil {
			result.Ok = true
			registeSuccNum++
		}

		str, _ := json.Marshal(&result)
		if err := util.CreateOrOpenFile(logFilename, string(str)); err != nil {
			s.log.Errorw("CreateOrOpenFile", err)
		}
		registeNum++

		if regError == nil {
			s.log.Infow("注册成功", strconv.Itoa(index+1)+" "+filename)
		} else {
			s.log.Infow("注册失败", strconv.Itoa(index+1)+" "+filename)
		}
	}

	s.log.Infow("【registeFaceOneByOne】end", "success", "注册的图片数量", registeNum, "注册成功的数量", registeSuccNum, "耗时", time.Since(t))
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
		if err := util.CreateOrOpenFile(logFilename, string(str)); err != nil {
			s.log.Errorw("CreateOrOpenFile", err)
		}

		if result.Ok {
			s.log.Infow("注册成功", strconv.Itoa(index+1)+" "+info.Filename)
		} else {
			s.log.Infow("注册失败", strconv.Itoa(index+1)+" "+info.Filename)
		}
	}

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
	if len(results) == 0 {
		err = ErrorFaceSearchEmpty
	}
	reply = &pb.SearchResultReply{}
	for _, result := range results {
		reply.Results = append(reply.Results, &pb.SearchResult{
			Filename: result.RegFilename,
			Match:    result.Match,
		})
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
