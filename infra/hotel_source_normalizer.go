package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-data-merge/usecase"
	"net/http"
	"strconv"
	"strings"
)

type HotelFetcher interface {
	GetHotels(ctx context.Context, httpClient *http.Client, name string) ([]usecase.Hotel, error)
}

type PaperfliesFetcher struct{}

func (n PaperfliesFetcher) GetHotels(ctx context.Context, httpClient *http.Client, name string) ([]usecase.Hotel, error) {
	endpoint := fmt.Sprintf("%s%s", url, name)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	var data []usecase.PaperfliesHotel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var hotels []usecase.Hotel
	for _, h := range data {
		hotel := normalizePaperfliesHotel(h)
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func normalizePaperfliesHotel(h usecase.PaperfliesHotel) usecase.Hotel {
	hotel := usecase.Hotel{
		HotelID:           h.HotelID,
		DestinationID:     h.DestinationID,
		Name:              h.HotelName,
		Description:       h.Details,
		BookingConditions: h.BookingConditions,
	}

	if h.Location != nil {
		hotel.Location = &usecase.HotelLocation{
			Address: h.Location.Address,
			Country: h.Location.Country,
		}
	}

	if h.Amenities != nil {
		hotel.Amenities = append(h.Amenities.GeneralAmenity, h.Amenities.RoomAmenity...)
	}

	if h.Images != nil {
		hotel.Images = &usecase.HotelImages{}

		for _, image := range h.Images.RoomImages {
			hotel.Images.RoomImages = append(hotel.Images.RoomImages, usecase.HotelImage{
				Link:        image.Link,
				Description: image.Caption,
			})
		}

		for _, image := range h.Images.SiteImages {
			hotel.Images.SiteImages = append(hotel.Images.SiteImages, usecase.HotelImage{
				Link:        image.Link,
				Description: image.Caption,
			})
		}
	}

	return hotel
}

type PatagoniaFetcher struct{}

func (n PatagoniaFetcher) GetHotels(ctx context.Context, httpClient *http.Client, name string) ([]usecase.Hotel, error) {
	endpoint := fmt.Sprintf("%s%s", url, name)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", endpoint, err)
	}

	var data []usecase.PatagoniaHotel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var hotels []usecase.Hotel
	for _, h := range data {
		hotel := normalizePatagoniaHotel(h)
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func normalizePatagoniaHotel(h usecase.PatagoniaHotel) usecase.Hotel {
	hotel := usecase.Hotel{
		HotelID:       h.HotelID,
		DestinationID: h.DestinationID,
		Name:          h.HotelName,
		Description:   h.Info,
		Location: &usecase.HotelLocation{
			Address:   h.Address,
			Latitude:  h.Latitude,
			Longitude: h.Longitude,
		},
		Amenities: h.Amenities,
	}

	if h.Images != nil {
		hotel.Images = &usecase.HotelImages{}

		for _, image := range h.Images.RoomImages {
			hotel.Images.RoomImages = append(hotel.Images.RoomImages, usecase.HotelImage{
				Link:        image.Url,
				Description: image.Description,
			})
		}

		for _, image := range h.Images.AmenityImages {
			hotel.Images.AmmenityImages = append(hotel.Images.AmmenityImages, usecase.HotelImage{
				Link:        image.Url,
				Description: image.Description,
			})
		}
	}

	return hotel
}

type AcmeFetcher struct{}

func (n AcmeFetcher) GetHotels(ctx context.Context, httpClient *http.Client, name string) ([]usecase.Hotel, error) {
	endpoint := fmt.Sprintf("%s%s", url, name)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	var data []usecase.AcmeHotel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var hotels []usecase.Hotel
	for _, h := range data {
		hotel := normalizeAcmeHotel(h)
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func normalizeAcmeHotel(h usecase.AcmeHotel) usecase.Hotel {
	address := h.Address

	if h.Address != nil && h.Postcode != nil && !strings.Contains(*h.Address, *h.Postcode) {
		add := fmt.Sprintf("%s, %s", strings.TrimSpace(*h.Address), strings.TrimSpace(*h.Postcode))
		address = &add
	}

	lat := float32Parser(h.Latitude)
	lng := float32Parser(h.Longitude)

	hotel := usecase.Hotel{
		HotelID:       h.HotelID,
		DestinationID: h.DestinationID,
		Name:          h.HotelName,
		Description:   h.Description,
		Location: &usecase.HotelLocation{
			Address:   address,
			City:      h.City,
			Country:   h.Country,
			Latitude:  lat,
			Longitude: lng,
		},
		Amenities: h.Facilities,
	}

	return hotel
}

func float32Parser(value interface{}) *float32 {
	switch v := value.(type) {
	case string:
		val, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil
		}
		float32Val := float32(val)
		return &float32Val
	case float32:
		return &v
	case float64:
		val := float32(v)
		return &val
	default:
		return nil
	}
}
