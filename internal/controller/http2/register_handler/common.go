package registerHandler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ParamsInt возвращает int значение соответствующего параметра из url.
// Если значение не строковое или его нет, вернется 0
func ParamsInt(ctx *fiber.Ctx, key string) int {
	str := ctx.Params(key)
	i, _ := strconv.Atoi(str)
	return i
}
