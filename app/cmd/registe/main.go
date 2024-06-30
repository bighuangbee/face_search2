package main

import (
	"flag"
	"fmt"
	"github.com/bighuangbee/face_search2/app/internal/service/face"
	"github.com/bighuangbee/face_search2/app/internal/service/face/face_recognize/face_wrapper"
	"github.com/bighuangbee/face_search2/pkg/conf"
	logger2 "github.com/bighuangbee/face_search2/pkg/logger"
	"github.com/bighuangbee/face_search2/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	"sync"
	"time"
)

var (
	flagconf  string
	logger    log.Logger
	bc        conf.Bootstrap
	regSerice *face.RegisteService
)

var regieteTime time.Duration

const FACE_REGISTE_FAILED = "/hiar_face/registe_failed"

func init() {
	flag.StringVar(&flagconf, "conf", "../../config", "config path, eg: -conf config.yaml")
	os.MkdirAll(FACE_REGISTE_FAILED, 0755)
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger = log.With(logger2.NewZapLogger(&logger2.Options{
		Level: zapcore.DebugLevel,
		Skip:  3,
		Writer: logger2.NewFileWriter(&logger2.FileOption{
			Filename: bc.Logger.Path + "/registe_%Y-%m-%d.log",
			MaxSize:  20,
		}),
	}))

	bc.Face.FaceMode = conf.FaceMode_registe
	if bc.Face.RegisteTimer <= 0 {
		bc.Face.RegisteTimer = 1
	}
	regieteTime = time.Minute * time.Duration(bc.Face.RegisteTimer)

	fmt.Println(1)
	//算法初始化
	app := face.NewFaceRecognizeApp(logger, &bc, nil)
	fmt.Println(2)
	regSerice = app.RegService

	fmt.Println(3)
	//处理队列，执行注册
	go regSerice.Run()
	fmt.Println(4)

	//定时检查新文件，放入注册队列
	go regieteHandle()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.POST("/unregiste", func(c *gin.Context) {
		logger.Log(log.LevelInfo, "【注册服务】UnRegisteAll", "注销照片")
		err := face_wrapper.UnRegisteAll()
		if err != nil {
			logger.Log(log.LevelError, "UnRegisteAll", err)
		}

	})

	router.POST("/registe", func(c *gin.Context) {
		readPhotosAndPushQueue()
	})

	logger.Log(log.LevelInfo, "照片注册服务启动", "")

	router.Run(":" + fmt.Sprintf("%d", 6666))
}

func readPhotosAndPushQueue() {
	regSerice.CheckExpired()

	registedSuccFace, registedFailedFace, newFace, err := regSerice.RegFilePreProcess()
	if err != nil && !os.IsNotExist(err) {
		logger.Log(log.LevelError, "注册文件预处理出错RegFilePreProcess", err)
		return
	}

	logger.Log(log.LevelInfo, "本次新增人脸", len(newFace), "之前注册成功人脸", len(registedSuccFace), "之前注册失败的人脸", len(registedFailedFace))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, filename := range newFace {
			regSerice.PushQueue(filename)
		}
	}()
	wg.Wait()

	// 通知搜索服务加载上一批数据
	_, port, _ := net.SplitHostPort(bc.Server.Http.Addr)
	resp, err := util.HttpPost(fmt.Sprintf("http://localhost:%s/face/reload", port), map[string]interface{}{})

	util.HttpPost(fmt.Sprintf("http://localhost:%s/face/reload", 6003), map[string]interface{}{})

	logger.Log(log.LevelInfo, "通知搜索服务加载数据", "", err, string(resp))
}

func regieteHandle() {
	readPhotosAndPushQueue()

	//logger.Log(log.LevelInfo, "测试注销  DSC00001", face.FACE_REGISTE_PATH+"/DSC00001.JPG")
	//if err := face_wrapper.UnRegiste(face.FACE_REGISTE_PATH + "/DSC00001.JPG"); err != nil {
	//	logger.Log(log.LevelError, "face_wrapper.UnRegiste error", err, "filename", "filename")
	//}

	regieteTimer := time.NewTicker(regieteTime)
	defer regieteTimer.Stop()

	for {
		<-regieteTimer.C
		logger.Log(log.LevelInfo, "定时执行注册照片", "")
		readPhotosAndPushQueue()
	}
}
