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
	Pool      *redis.Pool
	StartDate time.Time
	EndDate   time.Time
)

type PPV struct {
	Location     string
	Date         string
	Availability int
}

type User struct {
	MySejahteraID string `json:"mysejahteraId" redis:"mysejahteraId"`
	FirstName     string `json:"firstName" redis:"firstName"`
	LastName      string `json:"lastName" redis:"lastName"`
	Address       string `json:"address" redis:"address"`
	Location      string `json:"location" redis:"location"`
	PhoneNumber   string `json:"phoneNumber" redis:"phoneNumber"`
	Date          string `json:"date" redis:"date"`
}

func init() {
    StartDate = time.Date(2021, time.May, 15, 0, 0, 0, 0, time.UTC)
    EndDate = time.Date(2021, time.June, 15, 0, 0, 0, 0, time.UTC)

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

		ppvs, err := GetLocation("PWTC")
		if err != nil {
			c.SendString(err.Error())
		}

		// Set Cache-control header to 1s
		c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=1"))

		return c.JSON(ppvs)
	})

	app.Post("/submit", func(c *fiber.Ctx) error {

		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return err
		}

		ppv, err := GetPPV(user.Location, user.Date)
		if err != nil {
			return c.SendString("Failed to get PPV info")
		}
		if ppv.Availability <= 0 {
			return c.SendString("No more availability")
		}

		err = InsertUser(ppv, user)
		if err != nil {
			return c.SendString(err.Error())
		}

		// Publish message to Message Queue Broker

		return c.JSON(user)
	})

	app.Get("/submit", func(c *fiber.Ctx) error {

        // This endpoint is only for load testing
        // record insertion.
        // This endpoint can be invoked with GET method
        // and suitable for load testing tools like hey

		user := User{
			MySejahteraID: "1000",
			Date:          "20210601",
			Location:      "PWTC",
		}

		ppv, err := GetPPV("PWTC", "20210601")
		if ppv.Availability <= 0 {
			return c.SendString("No more availability")
		}

		err = InsertUser(ppv, &user)
		if err != nil {
			return c.SendString(err.Error())
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
	days := EndDate.Sub(StartDate).Hours() / 24
	daysInt := int(days)
	for i := 1; i < daysInt; i++ {
		date := StartDate.Add(time.Hour * time.Duration(i) * time.Duration(24))
		dateString := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
		ppv := PPV{Location: "PWTC", Date: dateString, Availability: 1000}
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

func GetPPV(location string, date string) (PPV, error) {
	conn := Pool.Get()
	defer conn.Close()

	// iterate each date
	key := "location:" + location + ":" + date

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
	ppv.Date = date
	ppv.Availability = availability

	return ppv, err
}

func GetLocation(location string) ([]PPV, error) {

	ppvs := []PPV{}

	// iterate each date
	days := EndDate.Sub(StartDate).Hours() / 24
	daysInt := int(days)
	for i := 1; i < daysInt; i++ {
		date := StartDate.Add(time.Hour * time.Duration(i) * time.Duration(24))
		dateString := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())

		ppv, err := GetPPV("PWTC", dateString)
		if err != nil {
			return ppvs, err
		}

		ppvs = append(ppvs, ppv)
	}
	return ppvs, nil
}

func InsertUser(ppv PPV, user *User) error {
	conn := Pool.Get()
	defer conn.Close()

	// Start a transaction
	conn.Send("MULTI")

	ppvKey := "location:" + ppv.Location + ":" + ppv.Date
	if err := conn.Send("DECR", ppvKey); err != nil {
		return fmt.Errorf("error setting key %s: %v", ppvKey, err)
	}

	userKey := "user:" + user.MySejahteraID
	if err := conn.Send("HMSET", redis.Args{}.Add(userKey).AddFlat(user)...); err != nil {
		return fmt.Errorf("error setting key %s: %v", userKey, err)
	}

	// Execute Transaction
	_, err := conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("error setting key %s: %v", userKey, err)
	}

	return nil
}

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
		Pool.Close()
		os.Exit(0)
	}()
}
