package face

import (
	"context"
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
	"os"
	"path/filepath"
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

const FACE_REGISTE_PATH = "/app/face_registe_path"

func NewFaceRecognizeApp(logger log.Logger, bc *conf.Bootstrap, data *data.Data) *FaceRecognizeApp {
	os.MkdirAll(FACE_REGISTE_PATH, 0755)

	face_models_path := os.Getenv("face_models_path")
	fmt.Println("face_models_path 1 ", face_models_path)
	if face_models_path != "" {
		face_models_path = "/root/face_search/libs/models/"
	}
	fmt.Println("face_models_path 2 ", face_models_path)

	app := FaceRecognizeApp{
		log:  log.NewHelper(log.With(logger, "module", "service/FaceRecognizeApp")),
		data: data,
	}

	err := face_wrapper.Init(face_models_path, "./hiarClusterLog.txt")
	app.log.Infow("【NewFaceRecognizeApp】face_wrapper init", err)

	return &app
}

func (s *FaceRecognizeApp) RegisteByPath(context.Context, *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	err := face_wrapper.UnRegisteAll()
	if err != nil {
		return nil, ErrorFaceSDK
	}

	go func() {
		t := time.Now()
		s.log.Infow("【RegisteByPath】begining", "")
		s.registering.Store(true)
		defer s.registering.Store(false)

		_, err := face_wrapper.Registe(FACE_REGISTE_PATH)
		if err != nil {
			s.log.Errorw("【face_wrapper.Registe】failed", err)
		} else {
			s.log.Infow("【RegisteByPath】end", "success", "duration", time.Since(t))
		}
	}()

	return &pb.EmptyReply{}, nil
}

func (s *FaceRecognizeApp) UnRegisteAll(ctx context.Context, req *pb.EmptyRequest) (*pb.EmptyReply, error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	err := face_wrapper.UnRegisteAll()
	if err != nil {
		return nil, ErrorFaceSDK
	}

	return &pb.EmptyReply{}, nil
}

func (s *FaceRecognizeApp) Search(ctx context.Context) (reply *pb.SearchResultReply, err error) {
	if s.registering.Load() {
		return nil, ErrorFaceRegistering
	}

	request, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return nil, ErrorRequestFrom
	}

	image, _, err := receiveFaceFile(request)
	if err != nil {
		return nil, err
	}

	results := face_wrapper.Search(image)
	if len(results) == 0 {
		return nil, ErrorFaceSearchEmpty
	}

	reply = &pb.SearchResultReply{}
	for _, result := range results {
		reply.Results = append(reply.Results, &pb.SearchResult{
			Filename: result.RegFilename,
			Match:    result.Match,
		})
	}

	return reply, nil
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
