package main

import "github.com/gofiber/fiber/v3"

func main() {
	app := fiber.New()

	app.Listen(":8080")
}
