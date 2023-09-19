package api

import (
	"fmt"

	"github.com/fulltimegodev/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore db.RoomStore
	
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore: rs,
	}
}

type HotelQueryParams struct {
	Rooms bool
	Raiting int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	
	var qparms HotelQueryParams
	
	if err := c.QueryParser(&qparms); err!= nil {
		return err
	}
	fmt.Println(qparms)
	
	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}