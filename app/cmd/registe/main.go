package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	v1 "github.com/bighuangbee/face_search2/api/biz/v1"
	"github.com/bighuangbee/face_search2/app/internal/service/face"
	"github.com/bighuangbee/face_search2/pkg/conf"
	logger2 "github.com/bighuangbee/face_search2/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"time"
)

var (
	flagconf string
	logger   log.Logger
	faceApp  *face.FaceRecognizeApp
	bc       conf.Bootstrap
)

var regieteTime = time.Minute * 3

func init() {
	flag.StringVar(&flagconf, "conf", "../../config", "config path, eg: -conf config.yaml")
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

	faceApp = face.NewFaceRecognizeApp(logger, &bc, nil)

	go regieteHandle()

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/unregiste", func(c *gin.Context) {
		//todo 删除db 重新注册
	})

	router.GET("/registe", func(c *gin.Context) {
		registe()
	})

	logger.Log(log.LevelInfo, "照片注册服务启动")

	router.Run(":" + fmt.Sprintf("%d", 6666))
}

func registe() {
	_, err := faceApp.RegisteByPath(context.Background(), &v1.RegisteRequest{Sync: true})
	if err != nil {
		logger.Log(log.LevelError, "注册照片，错误", err)
	} else {
		logger.Log(log.LevelInfo, "注册照片成功")

		//todo 通知搜索服务加载数据

		_, port, err := net.SplitHostPort(bc.Server.Http.Addr)
		url := fmt.Sprintf("http://localhost:%s/face/reload", port)

		data := map[string]interface{}{}
		jsonData, _ := json.Marshal(data)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Log(log.LevelInfo, "Failed to create request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Log(log.LevelInfo, "Failed to send reques:", err)
			return
		}
		defer resp.Body.Close()

		logger.Log(log.LevelInfo, "通知搜索服务加载数据", "")

	}
}

func regieteHandle() {

	registe()

	regieteTimer := time.NewTicker(regieteTime)
	defer regieteTimer.Stop()

	for {
		<-regieteTimer.C
		logger.Log(log.LevelInfo, "定时执行注册照片", "")
		registe()
	}
}
