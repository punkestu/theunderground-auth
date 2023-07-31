package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/handler"
	"github.com/punkestu/theunderground-auth/internal/lib"
	"github.com/punkestu/theunderground-auth/internal/repo/db"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fiber.New()

	conn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/theunderground?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	jwt := lib.NewJWT()
	r := db.NewDB(conn)
	e := entity.NewUserEntity(r)
	u := usecase.NewUserUsecase(e)
	h := handler.NewUserHandler(u, jwt)

	app.Post("/login", h.Login)
	app.Post("/key", h.LoginWithKey)
	app.Post("/register", h.Register)
	app.Post("/me", h.GetUser)

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			panic("ON APP LISTENING: " + err.Error())
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
mainLoop:
	for {
		switch <-sigChan {
		case syscall.SIGINT:
			break mainLoop
		}
	}
}
