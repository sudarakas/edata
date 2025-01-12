package subscription

import (
	"fmt"
	"net/http"
	"strconv"

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
	router.HandleFunc("/subscription/{code}", h.GetSubscriptionByCode).Methods("GET")
	router.HandleFunc("/subscription", h.GetAllSubscriptions).Methods("GET")
	router.HandleFunc("/subscription", h.UpdateSubscription).Methods("PUT")
	router.HandleFunc("/subscription", h.DeleteSubscription).Methods("DELETE")
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

// Get subscription by code
func (h *Handler) GetSubscriptionByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		// If the code is missing, return an error
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing code in request"))
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Get the subscription
	subscription, err := h.store.GetSubscriptionByCode(ctx, code)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the subscription
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "subscription retrieved successfully",
		"data": map[string]interface{}{
			"id":          subscription.SubscriptionID,
			"code":        subscription.Code,
			"plan":        subscription.Plan,
			"description": subscription.Description,
			"price":       subscription.Price,
			"data_limit":  subscription.DataLimit,
			"validity":    subscription.Validity,
		},
	})
}

// Get all subscriptions
func (h *Handler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	sort := r.URL.Query().Get("sort")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Convert limit and offset to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Get total count of subscriptions
	count, err := h.store.GetAllSubscriptionsCount(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch total count: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch subscriptions
	subscriptions, err := h.store.GetAllSubscriptions(r.Context(), sort, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	// Return paginated response
	utils.WritePaginatedSuccess(w, http.StatusOK, subscriptions, count, limit, offset)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		// If the code is missing, return an error
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing code in request"))
		return
	}

	// Parse the request payload
	var payload types.UpdateSubscriptionPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Check if the subscription exists
	existingSubscription, err := h.store.GetSubscriptionByCode(ctx, code)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Create a new subscription object, using payload values if provided,
	// otherwise falling back to the existing subscription values.
	subscription := types.Subscription{
		SubscriptionID: existingSubscription.SubscriptionID,
		Code:           payload.Code,
		Plan:           utils.GetOrDefault(payload.Plan, existingSubscription.Plan),
		Description:    utils.GetOrDefault(payload.Description, existingSubscription.Description),
		Price:          utils.GetOrDefaultFloat(payload.Price, existingSubscription.Price),
		DataLimit:      utils.GetOrDefaultInt(payload.DataLimit, existingSubscription.DataLimit),
		Validity:       utils.GetOrDefaultInt(payload.Validity, existingSubscription.Validity),
	}

	// Store the subscription
	newSubscription, err := h.store.UpdateSubscription(ctx, subscription)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the created subscription
	utils.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "subscription updated successfully",
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

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		// If the code is missing, return an error
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing code in request"))
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Check if the subscription exists
	existingSubscription, err := h.store.GetSubscriptionByCode(ctx, code)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Delete the subscription
	err = h.store.DeleteSubscription(ctx, existingSubscription.SubscriptionID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the created subscription
	utils.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "subscription deleted successfully",
	})
}
