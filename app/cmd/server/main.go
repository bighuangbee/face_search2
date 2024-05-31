package main

import (
	"flag"
	"github.com/bighuangbee/face_search2/pkg/conf"
	logger2 "github.com/bighuangbee/face_search2/pkg/logger"
	"github.com/bighuangbee/face_search2/pkg/util/kitGrpc"
	"github.com/go-kratos/kratos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"go.uber.org/zap/zapcore"
	"os"

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

func newApp(logger log.Logger, hs *http.Server, namingCli naming_client.INamingClient) *kratos.App {
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

	var naming naming_client.INamingClient
	if bc.Discovery.Enable {
		var err error
		naming, err = kitGrpc.NewNamingClient(
			bc.Discovery.Addr,
			uint64(bc.Discovery.Port),
			kitGrpc.WithLogLevel("info"),
			kitGrpc.WithNamespaceID(bc.Name),
		)
		if err != nil {
			panic(err)
		}
	}

	app, cleanup, err := wireApp(&bc, logger, naming)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
