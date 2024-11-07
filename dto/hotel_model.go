package dto

type ListHotelsRequest struct {
	HotelIDs       []string
	DestinationIDs []string
}

type ListHotelsResponse struct {
	Data []Hotel `json:"data"`
}

type Hotel struct {
	HotelID           string         `json:"hotel_id"`
	DestinationID     int32          `json:"destination_id"`
	Name              string         `json:"name"`
	Location          *HotelLocation `json:"location"`
	Description       string         `json:"description"`
	Amenities         *HotelAmenity  `json:"amenities"`
	Images            *HotelImages   `json:"images"`
	BookingConditions []string       `json:"booking_conditions"`
}

type HotelImages struct {
	RoomImages     []HotelImage `json:"rooms"`
	SiteImages     []HotelImage `json:"site"`
	AmmenityImages []HotelImage `json:"amenities"`
}

type HotelImage struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}

type HotelLocation struct {
	Latitude  *float32 `json:"latitude"`
	Longitude *float32 `json:"longitude"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	Country   *string  `json:"country"`
}

type HotelAmenity struct {
	GeneralAmenity []string `json:"general"`
	RoomAmenity    []string `json:"room"`
}
