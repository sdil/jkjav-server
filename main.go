package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/sdil/jkjav-server/pkg/station"
	"github.com/sdil/jkjav-server/api/routes"
	"github.com/sdil/jkjav-server/docs"
)

var (
	Pool      *redis.Pool
	StartDate time.Time
	EndDate   time.Time
)

func init() {
	// Redis connection establishment
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisPassword := os.Getenv("REDISPASSWORD")
	if redisHost == "" {
		redisHost = ""
		redisPort = "6379"
		redisPassword = ""
	}
	Pool = newPool(redisHost+":"+redisPort, redisPassword)

	// InitializeLocations()
	CleanupHook()
}

// @title JKJAV API Server
// @version 1.0
// @description High performant JKJAV API Server
// @BasePath /
func main() {

	RAILWAY_HOST := os.Getenv("RAILWAY_STATIC_URL")
	if RAILWAY_HOST == "" {
		docs.SwaggerInfo.Host = "localhost:3000"
	} else {
		docs.SwaggerInfo.Host = RAILWAY_HOST
	}

	stationRepo := station.NewRepo(Pool)
	stationService := station.NewService(stationRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/swagger/*", swagger.Handler)
	routes.StationRouter(app, stationService)
	// app.Post("/booking", routes.CreateBooking)

	app.Listen(":" + port)
}

// func InitializeLocations() {
// 	// Read locations.json file
// 	// And ensure the locations key in Redis exists
// 	log.Println("initialize locations")

// 	// iterate each date
// 	days := EndDate.Sub(StartDate).Hours() / 24
// 	daysInt := int(days)
// 	for i := 1; i < daysInt; i++ {
// 		date := StartDate.Add(time.Hour * time.Duration(i) * time.Duration(24))
// 		dateString := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
// 		ppv := PPV{Location: "PWTC", Date: dateString, Availability: 1000}
// 		SetLocation(ppv)
// 	}
// }

func newPool(server string, password string) *redis.Pool {

	log.Println("Connecting to redis server " + server)

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func CleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		fmt.Println("Shutting down!")
		Pool.Close()
		os.Exit(0)
	}()
}
