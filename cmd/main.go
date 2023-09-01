package main

import (
	"github.com/ShyunnY/cruise/logic"
	"github.com/ShyunnY/cruise/pkg/server"
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

	// build application for config
	app := ApplicationConfig()

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

func ApplicationConfig() *fiber.App {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     true,
		StrictRouting:         true,
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
	apiG.Get("/search", logic.SearchTracesService)
	apiG.Get("/query/:traceID", logic.QueryTracesService)
	apiG.Get("/service/:service", logic.ListTraceSvcService)
	apiG.Get("/operation/:operation", logic.ListOperationsService)

	return app
}
