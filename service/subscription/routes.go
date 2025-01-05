package subscription

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sudarakas/edata/types"
	"github.com/sudarakas/edata/utils"
)

type Handler struct {
	store types.SubscriptionStore
}

func NewHandler(store types.SubscriptionStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscription", h.CreateSubscription).Methods("POST")
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// Parse the request payload
	var payload types.CreateSubscriptionPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Valid the payload data
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Create a new subscription
	subscription := types.Subscription{
		Code:        payload.Code,
		Plan:        payload.Plan,
		Description: payload.Description,
		Price:       payload.Price,
		DataLimit:   payload.DataLimit,
		Validity:    payload.Validity,
	}

	// Store the subscription
	newSubscription, err := h.store.CreateSubscription(ctx, subscription)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the created subscription
	utils.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "subscription created successfully",
		"data": map[string]interface{}{
			"id":          newSubscription.SubscriptionID,
			"code":        newSubscription.Code,
			"plan":        newSubscription.Plan,
			"description": newSubscription.Description,
			"price":       newSubscription.Price,
			"data_limit":  newSubscription.DataLimit,
			"validity":    newSubscription.Validity,
		},
	})
}
