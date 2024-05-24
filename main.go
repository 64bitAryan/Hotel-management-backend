package main

import (
	"context"
	"flag"
	"log"

	"github.com/64bitAryan/hotel-management/api"
	"github.com/64bitAryan/hotel-management/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	port := flag.String("port", ":7545", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	var (
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		userStore   = db.NewMongoUserStore(client)
		userHandler = api.NewUserHandler(userStore)
		store       = db.Store{
			Room:  roomStore,
			Hotel: hotelStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(&store)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
	)
	// User haldlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// Hotel haldlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*port)
}
