package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewPostProfile(r service.UserReader) *postProfileHandler {
	return &postProfileHandler{
		reader: r,
	}
}

type postProfileHandler struct {
	reader service.UserReader
}

func (h *postProfileHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		//TODO do like this
		//var u service.UserLogin
		//err := ctx.BodyParser(&u)

		//TODO
		//1. Получаем токен и убеждаемся, что данные меняются для текущего пользователя
		//2. Обновляем данные
		//3. Передаем данные в форму через fiber.Map{}

		return ctx.Render("profile", fiber.Map{})
	}
}
