package main

import (
	"fmt"

	"github.com/64bitAryan/hotel-management/types"
)

func main() {
	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "France",
	}
	room := types.Room{
		Type:      types.SingleRoomType,
		BasePrice: 99.9,
	}
	fmt.Printf("%+v : %+v", hotel, room)
}
