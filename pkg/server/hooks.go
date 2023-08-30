package server

import "github.com/gofiber/fiber/v2"

func SetServerHooks(hooks *fiber.Hooks) {

	hooks.OnListen(ListenHook)
	hooks.OnShutdown()
}

func ListenHook(data fiber.ListenData) error {
	// TODO: 用于输出相关信息
	return nil
}

// ShutdownHook
// 用于处理优雅停机
func ShutdownHook() error {
	// TODO: 用于关闭相关组件
	return nil
}
