package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

type App struct {
	app *fiber.App
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("env var $PORT must be set")
	}

	a := App{app: fiber.New()}
	a.Run(port)
}

func (a *App) Run(port string) {
	a.app.Get("/all", func(c *fiber.Ctx) error {
		return c.JSON(GetAll())
	})
	a.app.Post("/add", func(c *fiber.Ctx) error {
		t := new(Task)
		if err := c.BodyParser(t); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if res, err := Add(*t); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.JSON(res)
		}
	})
	a.app.Put("/update", func(c *fiber.Ctx) error {
		ct := new(CombinedTask)
		if err := c.BodyParser(ct); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if res, err := Update(ct.ID.Id, ct.Task); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.JSON(res)
		}
	})
	a.app.Delete("/delete", func(c *fiber.Ctx) error {
		id := new(ID)
		if err := c.BodyParser(id); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if res, err := Delete(id.Id); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.JSON(res)
		}
	})

	err := a.app.Listen(":" + port)
	if err != nil {
		panic(err)
	}
}
