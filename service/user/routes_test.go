package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sudarakas/edata/types"
)

func TestUserServiceHandlers(t *testing.T) {
	// Create a mock store and handler
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayLoad{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "", // Invalid Email
			Password:  "XXXXXXXXXXX",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should register user correctly", func(t *testing.T) {
		payload := types.RegisterUserPayLoad{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "XXXXXXXXXXX",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) CreateUser(ctx context.Context, user types.User) (types.User, error) {
	return types.User{}, nil
}

func (m *mockUserStore) GetUserByEmail(ctx context.Context, email string) (types.User, error) {
	return types.User{}, nil
}

func (m *mockUserStore) GetUserByID(ctx context.Context, id string) (types.User, error) {
	return types.User{}, nil
}
