package main

import (
	"context"
	"flag"
	"log"

	"github.com/64bitAryan/hotel-management/api"
	"github.com/64bitAryan/hotel-management/api/middleware"
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
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = db.Store{
			Room:    roomStore,
			Hotel:   hotelStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		roomhandler  = api.NewRoomHandler(&store)
		hotelHandler = api.NewHotelHandler(&store)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		auth         = app.Group("/api")
	)

	//Auth handlers
	auth.Post("/auth", authHandler.HandleAuthentication)

	/*  Versioned API routes  */

	// User haldlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// Hotel haldlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Get("/room/", roomhandler.HandleGetRoom)
	apiv1.Post("/room/:id/book", roomhandler.HandleBookRoom)
	app.Listen(*port)
}
