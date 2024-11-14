package usecase

import (
	"context"
	"fmt"
	"hotel-data-merge/dto"
	"hotel-data-merge/pkg/cache"
	"strings"
	"time"
)

type HotelRepository interface {
	ListHotels(ctx context.Context) map[string][]Hotel
}

type HotelUsecase struct {
	hotelRepo HotelRepository
	cache     cache.CacheInterface
}

func NewHotelUsecase(repo HotelRepository, cache cache.CacheInterface) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo: repo,
		cache:     cache,
	}
}

const (
	GroupByDestination = "destination"
	GroupByHotel       = "hotel"
	CacheKey           = "hotels-cache-key"
)

func (u *HotelUsecase) ListHotels(ctx context.Context, req *dto.ListHotelsRequest) *dto.ListHotelsResponse {
	var mergedHotels map[string]Hotel
	var filteredIds []string
	var hotelsFromExternal map[string][]Hotel

	// filterType := GroupByDestination
	filteredIds = req.DestinationIDs

	// hotel id takes precedence for filter because it is more specific
	if len(req.HotelIDs) > 0 {
		filteredIds = req.HotelIDs
		// filterType = GroupByHotel
	}

	cacheVal, ok := u.cache.Get(CacheKey)

	if ok {
		hotelsFromExternal = cacheVal.(map[string][]Hotel)
	} else {
		hotelsFromExternal = u.hotelRepo.ListHotels(ctx)

		// this highly depends on how often the data changes
		u.cache.Set(CacheKey, hotelsFromExternal, 60*time.Minute)
	}

	mergedHotels = mergeHotelByID(hotelsFromExternal)

	hotelPartition := hotelPartitioning(mergedHotels)

	// return all hotels if there is no filter
	// filteredHotels := filterHotels(filterType, filteredIds, mergedHotels)
	filteredHotels := filterHotelsV2(filteredIds, hotelPartition)

	// add pagination here. page and limit
	cleanedHotels := cleanMergedData(filteredHotels)

	return &dto.ListHotelsResponse{
		Data: cleanedHotels,
	}
}

// map[string]Hotel -> map of the different id and the hotel detail
func hotelPartitioning(hotels map[string]Hotel) map[string]map[string][]Hotel {
	hotelPartition := map[string]map[string][]Hotel{}

	for _, hotel := range hotels {
		partitionKeyHotelID := string(hotel.HotelID[0])
		partitionKeyDestinationID := string(fmt.Sprintf("%d", hotel.DestinationID)[0])

		if _, exists := hotelPartition[partitionKeyHotelID]; exists {
			hotelPartition[partitionKeyHotelID][hotel.HotelID] = append(hotelPartition[partitionKeyHotelID][hotel.HotelID], hotel)
		} else {
			hotelPartition[partitionKeyHotelID] = map[string][]Hotel{
				hotel.HotelID: {hotel},
			}
		}

		if _, exists := hotelPartition[partitionKeyDestinationID]; exists {
			hotelPartition[partitionKeyDestinationID][fmt.Sprintf("%d", hotel.DestinationID)] = append(hotelPartition[partitionKeyDestinationID][fmt.Sprintf("%d", hotel.DestinationID)], hotel)
		} else {
			hotelPartition[partitionKeyDestinationID] = map[string][]Hotel{
				fmt.Sprintf("%d", hotel.DestinationID): {hotel},
			}
		}
	}

	return hotelPartition
}

func filterHotelsV2(ids []string, hotelPartioning map[string]map[string][]Hotel) map[string]Hotel {
	idsToInclude := map[string]bool{}
	for _, id := range ids {
		idsToInclude[id] = true
	}

	filteredHotels := map[string]Hotel{}

	for _, id := range ids {
		partitionKey := string(id[0])
		hotelsInPartition := hotelPartioning[partitionKey]
		fmt.Println(hotelsInPartition)
		// filteredHotels[id] = hotelsInPartition[id]
	}

	return filteredHotels
}

func filterHotels(filterType string, ids []string, mergedHotels map[string]Hotel) map[string]Hotel {
	if len(ids) == 0 {
		return mergedHotels
	}

	idsToInclude := map[string]bool{}
	for _, id := range ids {
		idsToInclude[id] = true
	}

	filteredHotels := map[string]Hotel{}

	if filterType == GroupByHotel {
		for id, hotelDetails := range mergedHotels {
			if _, exists := idsToInclude[id]; exists {
				filteredHotels[id] = hotelDetails
			}
		}

		return filteredHotels
	}

	for id, hotelDetails := range mergedHotels {
		descIDstr := fmt.Sprintf("%d", hotelDetails.DestinationID)
		if _, exists := idsToInclude[descIDstr]; exists {
			filteredHotels[id] = hotelDetails
		}
	}

	return filteredHotels
}

// cleanMergedData cleans the data and presents it in the api format we want to return
// cleaning includes trimming space and transforming data to returned format
func cleanMergedData(hotels map[string]Hotel) []dto.Hotel {
	cleanedHotels := []dto.Hotel{}

	for _, hotel := range hotels {
		cleanedHotel := dto.Hotel{
			HotelID:       hotel.HotelID,
			DestinationID: hotel.DestinationID,
			Name:          strings.TrimSpace(hotel.Name),
			Description:   strings.TrimSpace(hotel.Description),
			Amenities:     groupAmenity(hotel.Amenities),
			Images:        groupImages(hotel.Images),
			Location:      hotel.Location.toDto(),
		}

		bookingConditions := []string{}
		for _, bc := range hotel.BookingConditions {
			bookingConditions = append(bookingConditions, strings.TrimSpace(bc))
		}
		cleanedHotel.BookingConditions = bookingConditions

		cleanedHotels = append(cleanedHotels, cleanedHotel)
	}

	return cleanedHotels
}

// mergeHotelByID merges all 3 hotel sources data and groups them by hotel id
func mergeHotelByID(sources map[string][]Hotel) map[string]Hotel {
	mergedHotels := make(map[string]Hotel)

	for _, source := range sources {
		for _, hotel := range source {
			id := hotel.HotelID
			existingHotel, exists := mergedHotels[id]

			if !exists {
				mergedHotels[id] = hotel
				continue
			}

			if existingHotel.Location.Address == nil && hotel.Location.Address != nil {
				existingHotel.Location.Address = hotel.Location.Address
			}

			if existingHotel.Location.Latitude == nil && hotel.Location.Latitude != nil {
				existingHotel.Location.Latitude = hotel.Location.Latitude
			}

			if existingHotel.Location.Longitude == nil && hotel.Location.Longitude != nil {
				existingHotel.Location.Longitude = hotel.Location.Longitude
			}

			if existingHotel.Location.Country == nil && hotel.Location.Country != nil {
				existingHotel.Location.Country = hotel.Location.Country
			}

			if existingHotel.Location.City == nil && hotel.Location.City != nil {
				existingHotel.Location.City = hotel.Location.City
			}

			existingHotel.BookingConditions = append(existingHotel.BookingConditions, hotel.BookingConditions...)

			// choosing name based on length. but we can implement other scoring systems such as relevancy scoring
			if len(hotel.Name) > len(existingHotel.Name) {
				existingHotel.Name = hotel.Name
			}

			// choosing description based on length. but we can implement other scoring systems such as sentiments and relevancy scoring
			if len(hotel.Description) > len(existingHotel.Description) {
				existingHotel.Description = hotel.Description
			}

			// combining all the amenities and images first, will do normalization and removing of duplicates later
			existingHotel.Amenities = append(existingHotel.Amenities, hotel.Amenities...)

			if existingHotel.Images == nil {
				existingHotel.Images = hotel.Images
			} else if hotel.Images != nil {
				existingHotel.Images.AmmenityImages = append(existingHotel.Images.AmmenityImages, hotel.Images.AmmenityImages...)
				existingHotel.Images.SiteImages = append(existingHotel.Images.SiteImages, hotel.Images.SiteImages...)
				existingHotel.Images.RoomImages = append(existingHotel.Images.RoomImages, hotel.Images.RoomImages...)
			}

			mergedHotels[id] = existingHotel
		}
	}

	return mergedHotels
}

// groupImages groups the images together and removes duplicate images based on link and caption
func groupImages(images *HotelImages) *dto.HotelImages {
	cleanedImages := &dto.HotelImages{}

	amenityImages := map[string]string{}
	for _, image := range images.AmmenityImages {
		amenityImages[image.Link] = image.Description
	}

	siteImages := map[string]string{}
	for _, image := range images.SiteImages {
		siteImages[image.Link] = image.Description
	}

	roomImages := map[string]string{}
	for _, image := range images.RoomImages {
		roomImages[image.Link] = image.Description
	}

	for link, desc := range amenityImages {
		cleanedImages.AmmenityImages = append(cleanedImages.AmmenityImages, dto.HotelImage{
			Link:        link,
			Description: desc,
		})
	}

	for link, desc := range siteImages {
		cleanedImages.SiteImages = append(cleanedImages.SiteImages, dto.HotelImage{
			Link:        link,
			Description: desc,
		})
	}

	for link, desc := range roomImages {
		cleanedImages.RoomImages = append(cleanedImages.RoomImages, dto.HotelImage{
			Link:        link,
			Description: desc,
		})
	}

	return cleanedImages
}

// groupAmenity normalizes all the ammenties by ensuring they return the same amenity name,
// cleans the data (eg. trimspace) and remove any duplicates
func groupAmenity(amenities []string) *dto.HotelAmenity {
	generalAmenity := map[string]string{
		"pool":            "outdoor pool",
		"indoor pool":     "indoor pool",
		"outdoor pool":    "outdoor pool",
		"businesscenter":  "business center",
		"business centre": "business center",
		"business center": "business center",
		"wifi":            "wifi",
		"drycleaning":     "dry cleaning",
		"dry cleaning":    "dry cleaning",
		"breakfast":       "breakfast",
		"childcare":       "childcare",
		"parking":         "parking",
		"concierge":       "concierge",
		"bar":             "bar",
	}

	roomAmenity := map[string]string{
		"aircon":         "air conditioning",
		"tv":             "tv",
		"coffee machine": "coffee machine",
		"kettle":         "kettle",
		"hair dryer":     "hair dryer",
		"iron":           "iron",
		"tub":            "bath tub",
		"bathtub":        "bath tub",
		"minibar":        "minibar",
	}

	general := map[string]bool{}
	room := map[string]bool{}

	for _, amenity := range amenities {
		amenity = strings.TrimSpace(strings.ToLower(amenity))

		if val, exists := generalAmenity[amenity]; exists {
			general[val] = true
			continue
		}

		if val, exists := roomAmenity[amenity]; exists {
			room[val] = true
			continue
		}

		// if the amenity does not exist, we can log it here and update the list of amenities accordingly
		fmt.Printf("amenity not in mapping: %s", amenity)
	}

	hotelAmenity := &dto.HotelAmenity{}

	for key := range general {
		hotelAmenity.GeneralAmenity = append(hotelAmenity.GeneralAmenity, key)
	}

	for key := range room {
		hotelAmenity.RoomAmenity = append(hotelAmenity.RoomAmenity, key)
	}

	return hotelAmenity
}

// cleanCountryName parses the country name and returns a consistent value
func cleanCountryName(country *string) *string {
	countryMap := map[string]string{
		"singapore": "Singapore",
		"japan":     "Japan",
		"sg":        "Singapore",
		"jp":        "Japan",
	}

	c := strings.TrimSpace(strings.ToLower(*country))
	cleanedCountry := countryMap[c]

	return &cleanedCountry
}
