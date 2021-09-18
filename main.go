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
	"github.com/sdil/jkjav-server/api/routes"
	"github.com/sdil/jkjav-server/docs"
	"github.com/sdil/jkjav-server/pkg/booking"
	"github.com/sdil/jkjav-server/pkg/entities"
	"github.com/sdil/jkjav-server/pkg/station"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	StartDate = time.Date(2021, time.May, 15, 0, 0, 0, 0, time.UTC)
	EndDate   = time.Date(2021, time.June, 15, 0, 0, 0, 0, time.UTC)
)

// @title JKJAV API Server
// @version 1.0
// @description High performant JKJAV API Server
// @BasePath /
func main() {

	Pool := newRedisPool()
	defer Pool.Close()
	
	CleanupHook()

	if RAILWAY_HOST := os.Getenv("RAILWAY_STATIC_URL"); RAILWAY_HOST == "" {
		docs.SwaggerInfo.Host = "localhost:3000"
	} else {
		docs.SwaggerInfo.Host = RAILWAY_HOST
	}

	messageProducer := newKafkaProducer([]string{"localhost:9092"})

	stationRepo := station.NewRepo(Pool)
	stationService := station.NewService(stationRepo)
	InitializeLocations(stationService)

	bookingRepo := booking.NewRepo(Pool, messageProducer)
	bookingService := booking.NewService(bookingRepo)

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
	routes.BookingRouter(app, bookingService)

	app.Listen(":" + port)
}

func newKafkaProducer(brokerList []string) sarama.SyncProducer {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func InitializeLocations(service station.Service) {
	// Read locations.json file
	// And ensure the locations key in Redis exists
	log.Println("Initialize locations")

	// Iterate each date and create station
	days := EndDate.Sub(StartDate).Hours() / 24
	daysInt := int(days)
	for i := 1; i < daysInt; i++ {
		date := StartDate.Add(time.Hour * time.Duration(i) * time.Duration(24))
		dateString := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
		station := entities.Station{Location: "PWTC", Date: dateString, Availability: 1000}

		go func() {
			_, err := service.InsertStation(&station)
			if err != nil {
				log.Printf("Failed to create station %v", station)
			}
		}()
	}
}

func newRedisPool() *redis.Pool {

	// Redis connection establishment
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisPassword := os.Getenv("REDISPASSWORD")
	if redisHost == "" {
		redisHost = ""
		redisPort = "6379"
		redisPassword = ""
	}

	server := fmt.Sprintf("%s:%s", redisHost, redisPort)

	log.Println("Connecting to redis server " + server)

	return &redis.Pool{

		Wait:        true,
		MaxActive:   20,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(redisPassword))
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
		os.Exit(0)
	}()
}
