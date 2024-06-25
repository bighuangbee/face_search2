package main

import (
	"flag"
	"github.com/bighuangbee/face_search2/pkg/conf"
	logger2 "github.com/bighuangbee/face_search2/pkg/logger"
	"github.com/go-kratos/kratos/v2"
	"go.uber.org/zap/zapcore"
	"os"
	"os/exec"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flag conf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../config", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			// gs,
			hs,
		),
	)
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
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger := log.With(logger2.NewZapLogger(&logger2.Options{
		Level: zapcore.DebugLevel,
		Skip:  3,
		Writer: logger2.NewFileWriter(&logger2.FileOption{
			Filename: bc.Logger.Path + "/%Y-%m-%d.log",
			MaxSize:  20,
		}),
	}))

	logger = log.With(logger,
		//"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"service.id", id,
		//"service.name", Name,
		//"service.version", Version,
		//"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
	)

	//服务模式：搜索
	bc.Face.FaceMode = conf.FaceMode_search

	//todo 判断显存大小
	go registeRun(logger)

	app, cleanup, err := wireApp(&bc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func registeRun(logger log.Logger) {
	//等待主进程启动
	time.Sleep(10 * time.Second)

	//cmd := exec.Command("go", "run", "../registe/main.go")
	//cmd := exec.Command("./registe-bin")
	cmd := exec.Command("./registe-bin", "-conf", "/app/conf/config.yaml")

	logger.Log(log.LevelInfo, "启动注册服务registe-bin", "")
	err := cmd.Run()
	if err != nil {
		logger.Log(log.LevelError, "cmd.Run() failed", err)
	}
}
