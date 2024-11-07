package usecase

import (
	"context"
	"fmt"
	"hotel-data-merge/dto"
	"hotel-data-merge/pkg/cache"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupHotelTest() (*MockHotelRepository, *cache.MockCacheInterface) {
	mockHotelRepo := &MockHotelRepository{}
	mockCacheInterface := &cache.MockCacheInterface{}

	return mockHotelRepo, mockCacheInterface
}

func TestListHotels(t *testing.T) {
	mockHotelId := "mock-hotel-id"
	mockDestinationId := int32(1)
	mockAddress := "mock-address"
	mockCountry := "SG"
	mockDesc := "mock-desc"
	mockLongerDesc := "mock-longer-desc"
	mockBookingConditions := []string{"mock-booking-conditions"}
	mockLink := "mock-link"
	mockImgDesc := "mock-img-desc"
	mockHotelName := "mock-name"
	mockLatitude := float32(1.1)
	mockLongitude := float32(1.1)
	mockCity := "mock-city"
	mockPostcode := "mock-postcode"

	add := fmt.Sprintf("%s, %s", mockAddress, mockPostcode)
	normalizedHotels := func() map[string][]Hotel {
		return map[string][]Hotel{
			Paperflies: {
				{
					HotelID:           mockHotelId,
					DestinationID:     mockDestinationId,
					Name:              mockHotelName,
					Description:       mockLongerDesc,
					BookingConditions: mockBookingConditions,
					Location: &HotelLocation{
						Address: &mockAddress,
						Country: &mockCountry,
					},
					Amenities: []string{"coffee machine", "bar ", "aircon", "tv"},
					Images: &HotelImages{
						RoomImages: []HotelImage{
							{
								Link:        mockLink,
								Description: mockImgDesc,
							},
						},
						SiteImages: []HotelImage{
							{
								Link:        mockLink,
								Description: mockImgDesc,
							},
						},
					},
				},
			},
			Acme: {
				{
					HotelID:       mockHotelId,
					DestinationID: mockDestinationId,
					Name:          mockHotelName,
					Description:   mockDesc,
					Location: &HotelLocation{
						Address:   &add,
						Country:   &mockCountry,
						City:      &mockCity,
						Latitude:  &mockLatitude,
						Longitude: &mockLongitude,
					},
					Amenities: []string{"hair dryer", "outdoor pool", "wifi ", "minibar", "tv"},
				},
			},
			Patagonia: {
				{
					HotelID:       mockHotelId,
					DestinationID: mockDestinationId,
					Name:          mockHotelName,
					Description:   mockDesc,
					Location: &HotelLocation{
						Address:   &mockAddress,
						Longitude: &mockLongitude,
						Latitude:  &mockLatitude,
					},
					Amenities: []string{" pool", "outdoor pool", "wifi ", "aircon", "tv"},
					Images: &HotelImages{
						RoomImages: []HotelImage{
							{
								Link:        mockLink,
								Description: mockImgDesc,
							},
						},
						AmmenityImages: []HotelImage{
							{
								Link:        mockLink,
								Description: mockImgDesc,
							},
						},
					},
				},
			},
		}
	}

	cleanedCountry := "Singapore"
	mockReturnedHotels := func() []dto.Hotel {
		return []dto.Hotel{
			{
				HotelID:           mockHotelId,
				DestinationID:     mockDestinationId,
				Name:              mockHotelName,
				Description:       mockLongerDesc,
				BookingConditions: mockBookingConditions,
				Location: &dto.HotelLocation{
					Address:   &mockAddress,
					Country:   &cleanedCountry,
					Longitude: &mockLongitude,
					Latitude:  &mockLatitude,
					City:      &mockCity,
				},
				Amenities: &dto.HotelAmenity{
					GeneralAmenity: []string{"outdoor pool", "wifi", "bar"},
					RoomAmenity:    []string{"coffee machine", "air conditioning", "tv", "hair dryer", "minibar"},
				},
				Images: &dto.HotelImages{
					RoomImages: []dto.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
					SiteImages: []dto.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
					AmmenityImages: []dto.HotelImage{
						{
							Link:        mockLink,
							Description: mockImgDesc,
						},
					},
				},
			},
		}

	}

	t.Run("should successfully list hotels without filter and cache", func(t *testing.T) {
		mockHotelRepo, mockCache := setupHotelTest()
		usecase := NewHotelUsecase(mockHotelRepo, mockCache)
		ctx := context.Background()
		mockCache.On("Get", CacheKey).Return(nil, false)
		mockCache.On("Set", CacheKey, normalizedHotels(), 60*time.Minute)
		mockHotelRepo.On("ListHotels", ctx).Return(normalizedHotels())

		hotels := usecase.ListHotels(ctx, &dto.ListHotelsRequest{})

		assert.NotEmpty(t, hotels)
		assert.ElementsMatch(t, mockReturnedHotels()[0].Amenities.GeneralAmenity, hotels.Data[0].Amenities.GeneralAmenity)
		assert.ElementsMatch(t, mockReturnedHotels()[0].Amenities.RoomAmenity, hotels.Data[0].Amenities.RoomAmenity)
		assert.Equal(t, mockReturnedHotels()[0].Description, hotels.Data[0].Description)
		assert.Equal(t, mockReturnedHotels()[0].BookingConditions, hotels.Data[0].BookingConditions)
		assert.Equal(t, mockReturnedHotels()[0].Location, hotels.Data[0].Location)
		assert.Equal(t, mockReturnedHotels()[0].Images, hotels.Data[0].Images)
		assert.Equal(t, mockReturnedHotels()[0].Name, hotels.Data[0].Name)
		assert.Equal(t, mockReturnedHotels()[0].HotelID, hotels.Data[0].HotelID)
		mockCache.AssertExpectations(t)
		mockHotelRepo.AssertExpectations(t)
	})

	t.Run("should successfully list hotels with cache", func(t *testing.T) {
		mockHotelRepo, mockCache := setupHotelTest()
		usecase := NewHotelUsecase(mockHotelRepo, mockCache)

		mockCache.On("Get", CacheKey).Return(normalizedHotels(), true)
		hotels := usecase.ListHotels(context.Background(), &dto.ListHotelsRequest{})

		assert.NotEmpty(t, hotels)
		assert.ElementsMatch(t, mockReturnedHotels()[0].Amenities.GeneralAmenity, hotels.Data[0].Amenities.GeneralAmenity)
		assert.ElementsMatch(t, mockReturnedHotels()[0].Amenities.RoomAmenity, hotels.Data[0].Amenities.RoomAmenity)
		assert.Equal(t, mockReturnedHotels()[0].Description, hotels.Data[0].Description)
		assert.Equal(t, mockReturnedHotels()[0].BookingConditions, hotels.Data[0].BookingConditions)
		assert.Equal(t, mockReturnedHotels()[0].Location, hotels.Data[0].Location)
		assert.Equal(t, mockReturnedHotels()[0].Images, hotels.Data[0].Images)
		assert.Equal(t, mockReturnedHotels()[0].Name, hotels.Data[0].Name)
		assert.Equal(t, mockReturnedHotels()[0].HotelID, hotels.Data[0].HotelID)
		mockCache.AssertExpectations(t)
		mockHotelRepo.AssertExpectations(t)
	})

	t.Run("should successfully list only hotels filtered by destination id", func(t *testing.T) {
		mockHotelRepo, mockCache := setupHotelTest()
		usecase := NewHotelUsecase(mockHotelRepo, mockCache)

		normalizedHotels := normalizedHotels()
		mockCache.On("Get", CacheKey).Return(normalizedHotels, true)
		hotels := usecase.ListHotels(context.Background(), &dto.ListHotelsRequest{
			DestinationIDs: []string{"2"},
		})

		assert.Len(t, hotels.Data, 0)
		mockCache.AssertExpectations(t)
		mockHotelRepo.AssertExpectations(t)
	})

	t.Run("should successfully list only hotels filtered by hotel id", func(t *testing.T) {
		mockHotelRepo, mockCache := setupHotelTest()
		usecase := NewHotelUsecase(mockHotelRepo, mockCache)

		normalizedHotels := normalizedHotels()
		mockCache.On("Get", CacheKey).Return(normalizedHotels, true)
		hotels := usecase.ListHotels(context.Background(), &dto.ListHotelsRequest{
			HotelIDs:       []string{mockHotelId},
			DestinationIDs: []string{"2"},
		})

		assert.Len(t, hotels.Data, 1)
		mockCache.AssertExpectations(t)
		mockHotelRepo.AssertExpectations(t)
	})
}
