package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-data-merge/usecase"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CustomRoundTripper struct {
	roundTripFunc func(req *http.Request) *http.Response
}

func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.roundTripFunc(req), nil
}

func newMockClient(responseFunc func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &CustomRoundTripper{roundTripFunc: responseFunc},
	}
}

func TestListHotels(t *testing.T) {
	mockHotelId := "mock-hotel-id"
	mockDestinationId := int32(1)
	mockAddress := "mock-address"
	mockCountry := "mock-country"
	mockDesc := "mock-desc"
	mockBookingConditions := []string{"mock-booking-conditions"}
	mockGeneralAmenities := []string{"mock-general-amenities"}
	mockRoomAmenities := []string{"mock-room-amenities"}
	mockLink := "mock-link"
	mockImgDesc := "mock-img-desc"
	mockHotelName := "mock-name"
	mockLatitude := float32(1.1)
	mockLongitude := float32(1.1)
	mockCity := "mock-city"
	mockPostcode := "mock-postcode"

	mockPaperfliesHotels := []usecase.PaperfliesHotel{
		{
			HotelID:       mockHotelId,
			DestinationID: mockDestinationId,
			HotelName:     mockHotelName,
			Location: &usecase.Location{
				Address: &mockAddress,
				Country: &mockCountry,
			},
			Details:           mockDesc,
			BookingConditions: mockBookingConditions,
			Amenities: &usecase.Amenities{
				GeneralAmenity: mockGeneralAmenities,
				RoomAmenity:    mockRoomAmenities,
			},
			Images: &usecase.PaperfliesImages{
				RoomImages: []usecase.PaperfliesImage{
					{
						Link:    mockLink,
						Caption: mockImgDesc,
					},
				},
				SiteImages: []usecase.PaperfliesImage{
					{
						Link:    mockLink,
						Caption: mockImgDesc,
					},
				},
			},
		},
	}

	mockPatagoniaHotels := []usecase.PatagoniaHotel{
		{
			HotelID:       mockHotelId,
			DestinationID: mockDestinationId,
			HotelName:     mockHotelName,
			Address:       &mockAddress,
			Latitude:      &mockLatitude,
			Longitude:     &mockLongitude,
			Info:          mockDesc,
			Amenities:     []string{"mock-general-amenities", "mock-room-amenities"},
			Images: &usecase.PatagoniaImages{
				RoomImages: []usecase.PatagoniaImage{
					{
						Url:         mockLink,
						Description: mockImgDesc,
					},
				},
				AmenityImages: []usecase.PatagoniaImage{
					{
						Url:         mockLink,
						Description: mockImgDesc,
					},
				},
			},
		},
	}

	mockAcmeHotels := []usecase.AcmeHotel{
		{
			HotelID:       mockHotelId,
			DestinationID: mockDestinationId,
			HotelName:     mockHotelName,
			Address:       &mockAddress,
			Latitude:      &mockLatitude,
			Longitude:     &mockLongitude,
			City:          &mockCity,
			Country:       &mockCountry,
			Postcode:      &mockPostcode,
			Description:   mockDesc,
			Facilities:    []string{"mock-general-amenities", "mock-room-amenities"},
		},
	}

	add := fmt.Sprintf("%s, %s", mockAddress, mockPostcode)
	normalizedHotels := map[string][]usecase.Hotel{
		usecase.Paperflies: {
			{
				HotelID:           mockHotelId,
				DestinationID:     mockDestinationId,
				Name:              mockHotelName,
				Description:       mockDesc,
				BookingConditions: mockBookingConditions,
				Location: &usecase.HotelLocation{
					Address: &mockAddress,
					Country: &mockCountry,
				},
				Amenities: []string{"mock-general-amenities", "mock-room-amenities"},
				Images: &usecase.HotelImages{
					RoomImages: []usecase.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
					SiteImages: []usecase.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
				},
			},
		},
		usecase.Acme: {
			{
				HotelID:       mockHotelId,
				DestinationID: mockDestinationId,
				Name:          mockHotelName,
				Description:   mockDesc,
				Location: &usecase.HotelLocation{
					Address:   &add,
					Country:   &mockCountry,
					City:      &mockCity,
					Latitude:  &mockLatitude,
					Longitude: &mockLongitude,
				},
				Amenities: []string{"mock-general-amenities", "mock-room-amenities"},
			},
		},
		usecase.Patagonia: {
			{
				HotelID:       mockHotelId,
				DestinationID: mockDestinationId,
				Name:          mockHotelName,
				Description:   mockDesc,
				Location: &usecase.HotelLocation{
					Address:   &mockAddress,
					Longitude: &mockLongitude,
					Latitude:  &mockLatitude,
				},
				Amenities: []string{"mock-general-amenities", "mock-room-amenities"},
				Images: &usecase.HotelImages{
					RoomImages: []usecase.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
					AmmenityImages: []usecase.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
				},
			},
		},
	}

	responseFunc := func(req *http.Request) *http.Response {
		switch req.URL.Path {
		case fmt.Sprintf("/suppliers/%s", usecase.Paperflies):
			resp, _ := json.Marshal(mockPaperfliesHotels)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(string(resp))),
				Header:     make(http.Header),
			}
		case fmt.Sprintf("/suppliers/%s", usecase.Acme):
			resp, _ := json.Marshal(mockAcmeHotels)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(string(resp))),
				Header:     make(http.Header),
			}
		case fmt.Sprintf("/suppliers/%s", usecase.Patagonia):
			resp, _ := json.Marshal(mockPatagoniaHotels)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(string(resp))),
				Header:     make(http.Header),
			}
		default:
			resp, _ := json.Marshal(mockPaperfliesHotels)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(string(resp))),
				Header:     make(http.Header),
			}
		}
	}

	mockClient := newMockClient(responseFunc)

	t.Run("should successfully list and normalize hotels", func(t *testing.T) {
		r := NewHotelRepo(mockClient)
		hotels := r.ListHotels(context.Background())

		assert.NotEmpty(t, hotels)
		assert.Equal(t, normalizedHotels, hotels)
	})
}
