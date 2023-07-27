package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/handler"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"log"
)

func main() {
	app := fiber.New()

	// r := repo.NewUserRepo(conn)
	u := usecase.NewUserUsecase(nil) // TODO change the nil to repo
	h := handler.NewUserHandler(u)

	app.Post("/login", h.Login)
	app.Post("/register", h.Register)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatalln("ON APP LISTENING:", err.Error())
	}
}
