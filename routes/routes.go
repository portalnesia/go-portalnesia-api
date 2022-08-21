package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"portalnesia.com/api/middleware"
	"portalnesia.com/api/response"
)

func SetupRouters() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, errors error) error {
			err := response.Server(errors.Error())

			if e, ok := errors.(*response.Error); ok {
				// Override status code if fiber.Error type
				err = e
			}
			return c.Status(err.Status).JSON(fiber.Map{"data": nil, "error": err})
		},
		AppName: "Portalnesia v1",
	})

	app.Use(middleware.Database())

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	if os.Getenv("NODE_ENV") == "development" {
		app.Use(logger.New())
	}

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(etag.New())

	app.Use(requestid.New())

	app.Use(middleware.Header())

	app.Use(middleware.Authorization(middleware.AuthorizationConfig{Disable: true}))

	app.Static("/", os.Getenv("NODEJS_PUBLIC_PATH"))

	/*app.Use(limiter.New(limiter.Config{
		Max:               900,
		Expiration:        15 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))*/

	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "Portalnesia Metrics Page",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "API Uptime"})
	})

	app.Use(middleware.Authorization(middleware.AuthorizationConfig{}))

	v1 := app.Group("/v1")

	SetUpNewsRouters(v1)
	SetUpUserRouters(v1)

	app.Use(func(c *fiber.Ctx) error {
		err := response.EndpointNotFound()
		return c.Status(err.Status).JSON(fiber.Map{"data": nil, "error": err})
	})

	return app
}
