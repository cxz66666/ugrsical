package ugrsicalsrv

import "github.com/gofiber/fiber/v2"

func PingEp(ctx *fiber.Ctx) error {
	return ctx.SendString("pong")
}
