// fiber server example
// Copyright (C) 2024  Example

package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HTTPResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type FormRequest struct {
	Req   string `json:"req"`
	Count int    `json:"count"`
}

type FormResponse struct {
	Test string `json:"test"`
}

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		response := HTTPResponse{Data: "pong", Error: nil}
		return c.JSON(response)
	})

	app.Post("/postTest", func(c *fiber.Ctx) error {
		requestBody := new(FormRequest)
		if err := c.BodyParser(requestBody); err != nil {
			response := HTTPResponse{Data: nil, Error: "wrongData"}
			return c.JSON(response)
		}
		response := HTTPResponse{
			Data:  FormResponse{Test: fmt.Sprintf("%d", requestBody.Count)},
			Error: nil,
		}
		return c.JSON(response)
	})

	app.Listen(":3000")
}
