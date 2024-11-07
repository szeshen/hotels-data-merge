package usecase

import (
	"hotel-data-merge/dto"
	"strings"
)

const (
	Paperflies = "paperflies"
	Acme       = "acme"
	Patagonia  = "patagonia"
)

type Hotel struct {
	HotelID           string
	DestinationID     int32
	Name              string
	Location          *HotelLocation
	Description       string
	Amenities         []string // we will combine all amenities here then split them to general and room when we do our transformation later
	Images            *HotelImages
	BookingConditions []string
}

type HotelImages struct {
	RoomImages     []HotelImage
	SiteImages     []HotelImage
	AmmenityImages []HotelImage
}

type HotelImage struct {
	Link        string
	Description string
}

type HotelLocation struct {
	Latitude  *float32
	Longitude *float32
	Address   *string
	City      *string
	Country   *string
}

type HotelAmenity struct {
	GeneralAmenity []string
	RoomAmenity    []string
}

func (h *HotelLocation) toDto() *dto.HotelLocation {
	if h == nil {
		return nil
	}

	return &dto.HotelLocation{
		Latitude:  h.Latitude,
		Longitude: h.Longitude,
		Address:   trimspace(h.Address),
		City:      trimspace(h.City),
		Country:   cleanCountryName(h.Country),
	}
}

func trimspace(val *string) *string {
	if val == nil {
		return nil
	}

	v := strings.TrimSpace(*val)
	return &v
}

type PaperfliesHotel struct {
	HotelID           string            `json:"hotel_id"`
	DestinationID     int32             `json:"destination_id"`
	HotelName         string            `json:"hotel_name"`
	Location          *Location         `json:"location,omitempty"`
	Details           string            `json:"details,omitempty"`
	Amenities         *Amenities        `json:"amenities,omitempty"`
	Images            *PaperfliesImages `json:"images,omitempty"`
	BookingConditions []string          `json:"booking_conditions,omitempty"`
}

type Location struct {
	Address *string `json:"address,omitempty"`
	Country *string `json:"country,omitempty"`
}

type Amenities struct {
	GeneralAmenity []string `json:"general,omitempty"`
	RoomAmenity    []string `json:"room,omitempty"`
}

type PaperfliesImages struct {
	RoomImages []PaperfliesImage `json:"rooms,omitempty"`
	SiteImages []PaperfliesImage `json:"site,omitempty"`
}

type PaperfliesImage struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

type PatagoniaHotel struct {
	HotelID       string           `json:"id"`
	DestinationID int32            `json:"destination"`
	HotelName     string           `json:"name"`
	Latitude      *float32         `json:"lat,omitempty"`
	Longitude     *float32         `json:"lng,omitempty"`
	Address       *string          `json:"address,omitempty"`
	Info          string           `json:"info,omitempty"`
	Amenities     []string         `json:"amenities,omitempty"`
	Images        *PatagoniaImages `json:"images,omitempty"`
}

type PatagoniaImages struct {
	RoomImages    []PatagoniaImage `json:"rooms,omitempty"`
	AmenityImages []PatagoniaImage `json:"amenities,omitempty"`
}

type PatagoniaImage struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type AcmeHotel struct {
	HotelID       string      `json:"Id"`
	DestinationID int32       `json:"DestinationId"`
	HotelName     string      `json:"Name"`
	Latitude      interface{} `json:"Latitude,omitempty"`
	Longitude     interface{} `json:"Longitude,omitempty"`
	Address       *string     `json:"Address,omitempty"`
	City          *string     `json:"City,omitempty"`
	Country       *string     `json:"Country,omitempty"`
	Postcode      *string     `json:"PostalCode,omitempty"`
	Description   string      `json:"Description,omitempty"`
	Facilities    []string    `json:"Facilities,omitempty"`
}
