package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/sdil/jkjav-server/docs"
)

var (
	Pool        *redis.Pool
	StartDate   time.Time
	EndDate     time.Time
)

type PPV struct {
	Location     string `json:"location" redis:"location" example:"PWTC"`
	Date         string `json:"date" redis:"date" example:"20210516"`
	Availability int    `json:"availability" redis:"availability" example:"10"`
}

type User struct {
	MySejahteraID string `json:"mysejahteraId" redis:"mysejahteraId" example:"900127015527"`
	FirstName     string `json:"firstName" redis:"firstName" example:"Fadhil"`
	LastName      string `json:"lastName" redis:"lastName" example:"Yaacob"`
	Address       string `json:"address" redis:"address" example:"Kuala Lumpur"`
	Location      string `json:"location" redis:"location" example:"PWTC"`
	PhoneNumber   string `json:"phoneNumber" redis:"phoneNumber" example:"0123456789"`
	Date          string `json:"date" redis:"date" example:"20210516"`
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app := Setup()

	app.Listen(":" + port)
}

func Setup() *fiber.App {
		app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/swagger/*", swagger.Handler)
	app.Get("/list-ppv", ListPPV)
	app.Post("/submit", Submit)
	return app
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

// ListPPV godoc
// @Summary List PPV
// @Description Get PPV slots by state
// @Accept  json
// @Produce  json
// @Param state query string false "list PPV by state. The only available option is PWTC"
// @Success 200 {array} PPV
// @Router /list-ppv [get]
func ListPPV(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Please select a state")
	}

	ppvs, err := GetLocation("PWTC")
	if err != nil {
		c.SendString(err.Error())
	}

	// Set Cache-control header to 1s
	c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=1"))

	return c.JSON(ppvs)
}

// Submit godoc
// @Summary Submit
// @Description Submit vaccine booking slot
// @Accept  json
// @Produce  json
// @Param user body User true "User info"
// @Success 200 {object} User
// @Router /submit [post]
func Submit(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	err := InsertUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Publish message to Message Queue Broker

	return c.JSON(user)
}

func SetLocation(ppv PPV) error {
	conn := Pool.Get()
	defer conn.Close()

	key := "location:" + ppv.Location + ":" + ppv.Date

	_, err := conn.Do("SET", key, ppv.Availability)
	if err != nil {
		log.Printf("error setting key %s: %v", key, err)
		return fmt.Errorf("error setting key %s: %v", key, err)
	}

	log.Println("Added location " + key)
	return err
}

func GetLocation(location string) ([]PPV, error) {

	ppvs := []PPV{}

	// Iterate each date
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

func InsertUser(user *User) error {
	conn := Pool.Get()
	defer conn.Close()

	ppvKey := "location:" + user.Location + ":" + user.Date

	ok, err := redis.Bool(conn.Do("EXISTS", ppvKey))
	if err != nil {
		return fmt.Errorf("failed to location key %s. error: %v", ppvKey, err)
	}
	if ok == false {
		return fmt.Errorf("ppv %s location & date combination is invalid", ppvKey)
	}

	log.Printf("adding new user %s and updating ppv %s availability", user.MySejahteraID, ppvKey)

	// Start a transaction
	if _, err = conn.Do("WATCH", ppvKey); err != nil {
		return fmt.Errorf("Failed to watch key %s: %v", ppvKey, err)
	}

	var availability []byte
	availability, err = redis.Bytes(conn.Do("GET", ppvKey))
	if err != nil {
		return fmt.Errorf("error getting key %s: %v", ppvKey, err)
	}

	log.Printf("%s availability is %s. attempting to add user %s", ppvKey, string(availability), user.MySejahteraID)

	if string(availability) == "0" {
		if _, err = conn.Do("UNWATCH"); err != nil {
			log.Printf("%s failed to unwatch", ppvKey)
			return fmt.Errorf("Failed to unwatch key %s: %v", ppvKey, err)
		}

		log.Printf("%s is fully booked", ppvKey)
		return fmt.Errorf("Sorry, ppv is fully booked")
	}

	conn.Send("MULTI")

	// Decrease the counter
	if err := conn.Send("DECR", ppvKey); err != nil {
		return fmt.Errorf("error setting key %s: %v", ppvKey, err)
	}

	// Add new user dict
	userKey := "user:" + user.MySejahteraID
	if err := conn.Send("HMSET", redis.Args{}.Add(userKey).AddFlat(user)...); err != nil {
		return fmt.Errorf("error setting key %s: %v", userKey, err)
	}

	// Execute Transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("error setting key %s: %v", userKey, err)
	}

	log.Printf("successfully added new user %s and updated ppv %s availability", user.MySejahteraID, ppvKey)

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
