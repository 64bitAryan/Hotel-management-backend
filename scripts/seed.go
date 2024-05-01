package main

import (
	"context"
	"fmt"
	"log"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/64bitAryan/hotel-management/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "France",
	}
	room := types.Room{
		Type:      types.SingleRoomType,
		BasePrice: 99.9,
	}
	insertedHotel, err := hotelStore.InsetHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	room.HotelId = insertedHotel.ID
	fmt.Println(insertedHotel)
}
