package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

type PPV struct {
	Location     string
	Availability int
}

type User struct {
	MySejahteraID string
	FirstName     string
	LastName      string
	Address       string
	Location      string
	Date          string
}

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
	cleanupHook()
}

func main() {

	InitializeLocations()

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/list-ppv", func(c *fiber.Ctx) error {
		availablePPV := []PPV{}
		ppv := PPV{Location: "PWTC", Availability: 1000}
		availablePPV = append(availablePPV, ppv)
		return c.JSON(availablePPV)
	})

	app.Get("/submit", func(c *fiber.Ctx) error {
		user := User{
            MySejahteraID: "850113021157",
            FirstName: "Ahmad",
            LastName: "Albab",
            Address: "Lot 8-B, Taman Kulai, 14339, Johor Bahru",
            Location: "PWTC",
            Date: "20210530",
        }
        
        err := SetUser(user.MySejahteraID, user)
        if err != nil {
            c.SendString(err.Error())
        }
		// Publish message to Message Queue Broker
		return c.JSON(user)
	})

	app.Listen(":3000")
}

func InitializeLocations() {
	// Read locations.json file
	// And ensure the locations key in Redis exists
	log.Println("initialize locations")
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func SetUser(key string, user User) error {
    conn := Pool.Get()
	defer conn.Close()

    if _, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(user)...); err != nil {
		return fmt.Errorf("error setting key %s: %v", key, err)
	}

    return nil
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
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

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
