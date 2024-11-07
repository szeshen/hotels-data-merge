package srv

import (
	"context"
	"encoding/json"
	"hotel-data-merge/dto"
	"hotel-data-merge/usecase"
	"net/http"
	"strings"
)

type HotelHandler struct {
	hotelUsecase *usecase.HotelUsecase
}

func NewHotelHandler(hotelUsecase *usecase.HotelUsecase) *HotelHandler {
	return &HotelHandler{hotelUsecase: hotelUsecase}
}

func (h *HotelHandler) ListHotelsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	req := &dto.ListHotelsRequest{}

	hotelIDsStr := r.URL.Query().Get("hotel_ids")
	if hotelIDsStr != "" {
		req.HotelIDs = strings.Split(hotelIDsStr, ",")
	}
	destinationIDsStr := r.URL.Query().Get("destination_ids")
	if destinationIDsStr != "" {
		req.DestinationIDs = strings.Split(destinationIDsStr, ",")
	}

	hotel := h.hotelUsecase.ListHotels(ctx, req)
	json.NewEncoder(w).Encode(&hotel)
}
