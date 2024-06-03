package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/64bitAryan/hotel-management/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book room in past")
	}
	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		Store: store,
	}
}

func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	rooms, err := h.Store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
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

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID.String()),
		})
	}

	booking := &types.Booking{
		UserID:    user.ID,
		RoomId:    roomID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPerson: params.NumPerson,
	}
	inserted, err := h.Store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	c.JSON(inserted)
	return nil
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.Store.Booking.GetBooking(ctx, filter)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}

/*
return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID.String()),
		})
*/
