package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		request string
		path    string

		statusCode int
		response   string
	}{
		{
			name:   "get_empty",
			method: fiber.MethodGet,
			path:   "/all",

			statusCode: fiber.StatusOK,
			response:   "{}",
		},
		{
			name:    "add_one",
			method:  fiber.MethodPost,
			request: "{\"title\" : \"created task\",\"description\": \"created desc\",\"Completed\": false,\"Favorites\" : false}",
			path:    "/add",

			statusCode: fiber.StatusOK,
			response:   "{\"1\":{\"title\":\"created task\",\"description\":\"created desc\"}}",
		},
		{
			name:    "update_added",
			method:  fiber.MethodPut,
			request: "{\"id\": {\"id\": 1},\"task\": {\"title\": \"updated task\",\"description\": \"updated desc\",\"completed\": false,\"favorites\": false}}",
			path:    "/update",

			statusCode: fiber.StatusOK,
			response:   "{\"1\":{\"title\":\"updated task\",\"description\":\"updated desc\"}}",
		},
		{
			name:    "get_updated",
			method:  fiber.MethodGet,
			request: "{\"id\": {\"id\": 1},\"task\": {\"title\": \"updated task\",\"description\": \"updated desc\",\"completed\": false,\"favorites\": false}",
			path:    "/all",

			statusCode: fiber.StatusOK,
			response:   "{\"1\":{\"title\":\"updated task\",\"description\":\"updated desc\"}}",
		},
		{
			name:    "delete_updated",
			method:  fiber.MethodDelete,
			request: "{\"id\": 1}",
			path:    "/delete",

			statusCode: fiber.StatusOK,
			response:   "{}",
		},
	}

	app := fiber.New()
	app.Get("/all", func(c *fiber.Ctx) error {
		return c.JSON(GetAll())
	})
	app.Post("/add", func(c *fiber.Ctx) error {
		b := new(Task)
		if err := c.BodyParser(b); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if res, err := Add(*b); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.JSON(res)
		}
	})
	app.Put("/update", func(c *fiber.Ctx) error {
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
	app.Delete("/delete", func(c *fiber.Ctx) error {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, bytes.NewBuffer([]byte(test.request)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			if resp.StatusCode != test.statusCode {
				fmt.Printf("expected status code: %d, actual %d", resp.StatusCode, test.statusCode)
				t.Fail()
			}

			if body, _ := ioutil.ReadAll(resp.Body); string(body) != test.response {
				fmt.Printf("expected reponse: %s, actual %s", test.response, string(body))
				t.Fail()
			}
		})
	}
}
