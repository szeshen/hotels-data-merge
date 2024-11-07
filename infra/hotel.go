package infra

import (
	"hotel-data-merge/usecase"
	"net/http"
	"sync"
)

const (
	url = "https://5f2be0b4ffc88500167b85a0.mockapi.io/suppliers/"
)

type HotelSourceConfig struct {
	name         string
	hotelFetcher HotelFetcher
}

type HotelRepo struct {
	httpClient         *http.Client
	hotelSourceConfigs []HotelSourceConfig
}

func NewHotelRepo(client *http.Client) usecase.HotelRepository {
	if client == nil {
		client = &http.Client{}
	}

	return &HotelRepo{
		httpClient: client,
		hotelSourceConfigs: []HotelSourceConfig{
			{
				name:         usecase.Patagonia,
				hotelFetcher: PatagoniaFetcher{},
			},
			{
				name:         usecase.Paperflies,
				hotelFetcher: PaperfliesFetcher{},
			},
			{
				name:         usecase.Acme,
				hotelFetcher: AcmeFetcher{},
			},
		},
	}
}

func (hr *HotelRepo) ListHotels() map[string][]usecase.Hotel {
	hotels := map[string][]usecase.Hotel{}
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, config := range hr.hotelSourceConfigs {
		wg.Add(1)
		go func(config HotelSourceConfig) {
			defer wg.Done()
			normalizedHotels, err := config.hotelFetcher.GetHotels(hr.httpClient, config.name)
			if err != nil {
				return
			}

			mutex.Lock()
			hotels[config.name] = normalizedHotels
			mutex.Unlock()
		}(config)
	}
	wg.Wait()

	return hotels
}
