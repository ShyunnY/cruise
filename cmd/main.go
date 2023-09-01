package main

import (
	"context"
	"github.com/ShyunnY/cruise/logic"
	"github.com/ShyunnY/cruise/pkg/metricx"
	pipe "github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/pool"
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/ShyunnY/cruise/pkg/server"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// TODO: 加入config配置
	// get config path by cli-flag
	//var configPath string
	//flag.StringVar(&configPath, "config.file", "", "The absolute path to the cruise configuration file.")
	//flag.Parse()
	//
	//// get config path by env
	//env := os.Getenv("CRUISE_CONFIG")
	//if env != "" {
	//	configPath = env
	//} else if configPath == "" {
	//	log.Println("The configuration path cannot be empty,the configuration must be provided")
	//	os.Exit(1)
	//}
	//
	//// read configuration file
	//confFile, err := os.ReadFile(filepath.Clean(configPath))
	//if err != nil {
	//	log.Printf("cannot read config file, config: %s err: %s", configPath, err)
	//	os.Exit(1)
	//}
	//
	//// unmarshal yaml file
	//yaml.Unmarshal(confFile, nil)

	// start component
	//RunComponent()

	// ⭐: build application for config
	app := ApplicationConfig(logic.ServiceCtx{})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("graceful shutdown")
	// graceful shutdown
	gracefulDuration := time.Second * 10
	if err := app.ShutdownWithTimeout(gracefulDuration); err != nil {
		log.Println(err)
	}
}

// RunComponent
// build criuse component by Config
func RunComponent(stopCh <-chan struct{}) error {

	statfn := func() error {

		// todo: 根据配置文件选择reader
		grpcReader, err := reader.NewGrpcReader(reader.GrpcReaderConfig{
			Host: "192.168.136.134",
			Port: 16685,
		})
		if err != nil {
			return err
		}

		// todo: 根据配置文件选择store
		storeMemory := memory.NewStoreMemory()

		// todo: 根据配置文件选择 sink
		sink := pipe.NewStorageSink(pipe.StorageSinkConfig{
			Store:     storeMemory,
			BatchSize: 4,
			Interval:  time.Second * 30,
		})

		manage := metricx.NewManage(metricx.ManageConfig{})
		metricxStage := pipe.NewMetricxStage(manage)

		// todo: 根据配置文件进行配置pool
		workPool := pool.NewWorkPool(pool.WorkConfig{
			Interval: time.Second * 20,
			Read:     grpcReader,
			Sink:     sink,
			Stages:   []pipe.StagePipe{metricxStage},
			BufSize:  2,
			Period:   time.Hour * 24,
		})
		workPool.Work(context.TODO())

		return nil
	}

	statfn()
	return nil
}

// ApplicationConfig
// build http application by config
func ApplicationConfig(svcCtx logic.ServiceCtx) *fiber.App {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     true,
		StrictRouting:         true,

		// TODO: 需要考虑如何对handler func产生的error如何处理
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {

			return nil
		},

		// TODO: 我们后续可以考虑是否自定义sonic的Config
		// sonic marchsal as default JSONEncoder
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return sonic.Marshal(v)
		},
		// sonic unmarshal as default JSONDecoder
		JSONDecoder: func(d []byte, v interface{}) error {
			return sonic.Unmarshal(d, v)
		},
	})

	// init fiber.hook
	server.SetServerHooks(app.Hooks())
	// init fiber.middleware
	server.SetServerMiddleware(app)

	// search trace read group
	// compatibility otel Read interface(jaeger_v3 interface)
	apiG := app.Group("/read")
	apiG.Get("/search", logic.SearchTracesService(svcCtx))
	apiG.Get("/query/:traceID", logic.QueryTracesService(svcCtx))
	apiG.Get("/service/:service", logic.ListTraceSvcService(svcCtx))
	apiG.Get("/operation/:operation", logic.ListOperationsService(svcCtx))

	return app
}
