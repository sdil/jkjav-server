package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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
	Date         string
	Availability int
}

type User struct {
	MySejahteraID string
	FirstName     string
	LastName      string
	Address       string
	Location      string
	PhoneNumber   string
	Date          string
}

func init() {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
	InitializeLocations()
	CleanupHook()
}

func main() {

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/list-ppv", func(c *fiber.Ctx) error {

        state := c.Query("state")
        if state == "" {
            return c.SendString("Please select a state")
        }
        
        availablePPV := []PPV{}
		ppv, err := GetLocation("PWTC")
		if err != nil {
            c.SendString(err.Error())
		}
        
		availablePPV = append(availablePPV, ppv)

        // Set Cache-control header to 1s
        c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=1"))

		return c.JSON(availablePPV)
	})

	app.Get("/submit", func(c *fiber.Ctx) error {

		user := User{
			MySejahteraID: "850113021157",
			FirstName:     "Ahmad",
			LastName:      "Albab",
			Address:       "Lot 8-B, Taman Kulai, 14339, Johor Bahru",
			Location:      "PWTC",
			PhoneNumber:   "0149221442",
			Date:          "20210530",
		}
		ppv := PPV{Location: "PWTC", Date: "2021-05-01"}

		err := SetUser(ppv, user)
		if err != nil {
			c.SendString(err.Error())
		}

		// Publish message to Message Queue Broker

		return c.JSON(user)
	})

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

	app.Listen(":" + port)
}

func InitializeLocations() {
	// Read locations.json file
	// And ensure the locations key in Redis exists
	log.Println("initialize locations")

	// iterate each date
	ppv1 := PPV{Location: "PWTC", Date: "2021-05-01", Availability: 1000}
	ppvs := []PPV{}
	ppvs = append(ppvs, ppv1)

	for _, ppv := range ppvs {
		SetLocation(ppv)
	}
}

func SetLocation(ppv PPV) error {
	conn := Pool.Get()
	defer conn.Close()

	key := "location:" + ppv.Location + ":" + ppv.Date

	_, err := conn.Do("SET", key, ppv.Availability)
	if err != nil {
        log.Println(fmt.Sprintf("error setting key %s: %v", key, err))
		return fmt.Errorf("error setting key %s: %v", key, err)
	}

	log.Println("Added location " + key)
	return err
}

func GetLocation(location string) (PPV, error) {

	conn := Pool.Get()
	defer conn.Close()

	// iterate each date
	key := "location:" + location + ":2021-05-01"

	ppv := PPV{}
	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return ppv, fmt.Errorf("error getting key %s: %v", key, err)
	}

	availability, err := strconv.Atoi(string(data))
	if err != nil {
		log.Println("error reading location availability: " + err.Error())
	}

	ppv.Location = location
	ppv.Date = "2021-05-01"
	ppv.Availability = availability

	return ppv, err
}

func SetUser(ppv PPV, user User) error {
	conn := Pool.Get()
	defer conn.Close()

	// Multi
	ppvKey := "location:" + ppv.Location + ":" + ppv.Date
	if _, err := redis.Int(conn.Do("DECR", ppvKey)); err != nil {
		return fmt.Errorf("error setting key %s: %v", ppvKey, err)
	}

	userKey := "user:" + user.MySejahteraID
	if _, err := conn.Do("HMSET", redis.Args{}.Add(userKey).AddFlat(user)...); err != nil {
		return fmt.Errorf("error setting key %s: %v", userKey, err)
	}
	// Commit

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

func CleanupHook() {

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
