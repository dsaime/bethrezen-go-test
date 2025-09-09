package registerHandler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ParamsInt(ctx *fiber.Ctx, key string) int {
	str := ctx.Params(key)
	i, _ := strconv.Atoi(str)
	return i
}
