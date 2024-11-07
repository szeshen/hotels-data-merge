package main

import (
	"hotel-data-merge/infra"
	"hotel-data-merge/pkg/cache"
	"hotel-data-merge/srv"
	"hotel-data-merge/usecase"
	"log"
	"net/http"
	"time"
)

func main() {
	repo := infra.NewHotelRepo(nil)
	cache := cache.NewGoCacheWrapper(60*time.Minute, 75*time.Minute)
	usecase := usecase.NewHotelUsecase(repo, cache)
	handler := srv.NewHotelHandler(usecase)

	// Set up HTTP server
	http.HandleFunc("/hotels", handler.ListHotelsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
