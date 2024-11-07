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
	Description       string         `json:"description,omitempty"`
	Amenities         *HotelAmenity  `json:"amenities,omitempty"`
	Images            *HotelImages   `json:"images,omitempty"`
	BookingConditions []string       `json:"booking_conditions,omitempty"`
}

type HotelImages struct {
	RoomImages     []HotelImage `json:"rooms,omitempty"`
	SiteImages     []HotelImage `json:"site,omitempty"`
	AmmenityImages []HotelImage `json:"amenities,omitempty"`
}

type HotelImage struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}

type HotelLocation struct {
	Latitude  *float32 `json:"latitude,omitempty"`
	Longitude *float32 `json:"longitude,omitempty"`
	Address   *string  `json:"address,omitempty"`
	City      *string  `json:"city,omitempty"`
	Country   *string  `json:"country,omitempty"`
}

type HotelAmenity struct {
	GeneralAmenity []string `json:"general,omitempty"`
	RoomAmenity    []string `json:"room,omitempty"`
}
