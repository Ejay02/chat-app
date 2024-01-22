// package main

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/gofiber/fiber/v3/middleware/cors"

// 	// "github.com/pusher/pusher-http-go"
// 	"github.com/pusher/pusher-http-go/v5"
// )

// func main() {
//     app := fiber.New()

// 		app.Use(cors.New())

// 		pusherClient := pusher.Client{
// 			AppID: "1511755",
// 			Key: "70c3c70144d1c8f33bdb",
// 			Secret: "2669dcfe8571069e8e1e",
// 			Cluster: "mt1",
// 			Secure: true,
// 		}

//     app.Post("/api/messages", func(c fiber.Ctx) error {
// 			data := map[string]string{"message": "hello world"}
// 			err := pusherClient.Trigger("chat", "message", data)
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			return c.JSON([]string{})
//     })

//     app.Listen(":8000")
// }

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Use(cors.New())

	pusherClient := pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
		Secure:  true,
				
	}

	app.Post("/api/messages", func(c *fiber.Ctx) error {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			fmt.Println(err.Error())
			return err
		}

		message, exists := data["message"]
		if !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Message field not found",
			})
		}

		err := pusherClient.Trigger("chat", "message", map[string]string{"message": message})
		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return c.JSON([]string{})
	})

	app.Listen(":8000")
}
