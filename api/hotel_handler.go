package api

import (
	"fmt"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	HotelStore db.HotelStore
	RoomStore  db.RoomStore
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		HotelStore: hs,
		RoomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	fmt.Println(qparams)
	res, err := h.HotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
