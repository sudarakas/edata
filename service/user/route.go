package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sudarakas/edata/service/auth"
	"github.com/sudarakas/edata/types"
	"github.com/sudarakas/edata/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoute(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.Handle("/user", auth.WithJWTAuth(http.HandlerFunc(h.handleGetUserByCondition))).Methods("GET")
	router.Handle("/user", auth.WithJWTAuth(http.HandlerFunc(h.handleUpdateUser))).Methods("PATCH")
	router.Handle("/user", auth.WithJWTAuth(http.HandlerFunc(h.handleDeleteUser))).Methods("DELETE")
	router.Handle("/user/change-password", auth.WithJWTAuth(http.HandlerFunc(h.handleGetUser))).Methods("PATCH")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Get the payload data
	var payload types.LoginUserPayLoad
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

	// Check if the user exists
	user, err := h.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Check if the password is correct
	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Generate the token
	token, err := auth.GenerateJWT(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success message
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "login successful",
		"token":   token,
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {

	// Get the payload data
	var payload types.RegisterUserPayLoad
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

	// Check if the user exists
	existingUser, _ := h.store.GetUserByEmail(ctx, payload.Email)
	if existingUser.ID != "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Create the new user
	user := types.User{
		Email:     payload.Email,
		Password:  hashedPassword,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the user
	newUser, err := h.store.CreateUser(ctx, user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success message
	utils.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
		"data": map[string]interface{}{
			"firstName": newUser.FirstName,
			"lastName":  newUser.LastName,
			"email":     newUser.Email,
		},
	})

}

func (h *Handler) handleGetUserByCondition(w http.ResponseWriter, r *http.Request) {
	// Check if email query parameter exists
	email := r.URL.Query().Get("email")

	if email != "" {
		// If email is provided, fetch user by email
		h.handleGetUserByEmail(w, r)
	} else {
		// Otherwise, fetch user by ID or default logic
		h.handleGetUser(w, r)
	}
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from the context
	claims, ok := r.Context().Value(auth.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userID, ok := claims["id"].(string) // Assuming ID is a string
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Retrieve the user using the ID
	user, err := h.store.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user details
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "User retrieved successfully",
		"user": map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

func (h *Handler) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	// Extract the email query parameter from the URL
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email query parameter is required", http.StatusBadRequest)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Retrieve the user using the email
	user, err := h.store.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user details
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "User retrieved successfully",
		"user": map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from the context
	claims, ok := r.Context().Value(auth.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Retrieve the user using the ID
	user, err := h.store.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get the payload data
	var payload types.UpdateUserPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Valid the payload data
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Update the user
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Email = payload.Email

	// Update the user
	updatedUser, err := h.store.UpdateUser(ctx, user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success message
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "user updated successfully",
		"user": map[string]interface{}{
			"firstName": updatedUser.FirstName,
			"lastName":  updatedUser.LastName,
			"email":     updatedUser.Email,
			"createdAt": updatedUser.CreatedAt,
			"updatedAt": updatedUser.UpdatedAt,
		},
	})
}

func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from the context
	claims, ok := r.Context().Value(auth.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Retrieve the user using the ID
	user, err := h.store.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete the user
	err = h.store.DeleteUser(ctx, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success message
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})
}

func (h *Handler) changeUserPassword(w http.ResponseWriter, r *http.Request) {
	// Get claims from the context
	claims, ok := r.Context().Value(auth.ClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Use context from the request
	ctx := r.Context()

	// Retrieve the user using the ID
	user, err := h.store.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get the payload data
	var payload types.ChangePasswordPayLoad
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Valid the payload data
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// Check if the old password is correct
	if !auth.CheckPasswordHash(payload.OldPassword, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid old password"))
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Update the user's password
	user.Password = hashedPassword

	// Update the user
	err = h.store.ChangePassword(ctx, user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return success message
	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "password changed successfully",
	})
}
