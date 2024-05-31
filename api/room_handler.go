package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/64bitAryan/hotel-management/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	Store *db.Store
}

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	TillDate  time.Time `json:"tillDate"`
	NumPerson int       `json:"numPerson"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		Store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	booking := types.Booking{
		UserID:    user.ID,
		RoomId:    roomID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPerson: params.NumPerson,
	}
	fmt.Println(booking)
	fmt.Printf("%+v\n", booking)
	return nil
}
