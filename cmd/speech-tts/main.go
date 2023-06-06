package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"speech-tts/internal/conf"
	"speech-tts/internal/pkg/log"
	"speech-tts/internal/pkg/nacos"
	"speech-tts/internal/pkg/trace"
	"speech-tts/internal/utils"
	"syscall"

	_ "go.uber.org/automaxprocs"
	_ "speech-tts/internal/pkg/catch"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	Commit  string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	utils.SetServerVersion(Version, Commit)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, config *conf.Data) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Signal(syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGSEGV),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(nacos.NewRegister(config)),
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

	myLogger := Mylog.NewLogger(bc.Log)
	myLogger.SetLogger()
	trace.InitTracer(bc.Data.Otel.Addr, bc.Data.Otel.Name)

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Data.App.Path, myLogger.GetLogger())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
