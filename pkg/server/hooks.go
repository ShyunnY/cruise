package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func SetServerHooks(hooks *fiber.Hooks) {
	hooks.OnListen(listenHook)
	hooks.OnShutdown(shutdownHook)
}

func listenHook(data fiber.ListenData) error {
	// 启动相关组件
	log.Printf("cruise complate running on %s:%s", data.Host, data.Port)
	return nil
}

// ShutdownHook
// 用于处理优雅停机
func shutdownHook() error {
	// TODO: 关闭相关组件

	return nil
}
