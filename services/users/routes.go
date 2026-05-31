package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YhomiAce/rest-api-pg/config"
	"github.com/YhomiAce/rest-api-pg/types"
	"github.com/YhomiAce/rest-api-pg/utils"
	"github.com/gorilla/mux"
)


type Handler struct {
	store types.UserStore
	config *config.Config
}

func NewHandler(store types.UserStore, config *config.Config) *Handler {
	return &Handler{store: store, config: config}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}
	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// check if user exists
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to check user existence: %v", err))
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	
	// compare password
	if !utils.ComparePassword(payload.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	// Parse the duration string
	expiresIn, err := time.ParseDuration(h.config.JWTExpiresIn)
	if err != nil {
		// Handle error - use default value or return error
		log.Printf("Invalid JWT expiry duration: %v, using default 1h", err)
		expiresIn = time.Hour
	}
	token, err := utils.GenerateJWT(user.ID, h.config.JWTSecret, expiresIn)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate JWT token: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user": user,
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// check if user already exists
	existingUser, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to check user existence: %v", err))
		return
	}
	if existingUser != nil {
		utils.WriteError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to hash password: %v", err))
		return
	}
	
	// create user in database (this part is simplified and should be implemented properly)
	// In a real application, you would also want to return the created user or a success message
	err = h.store.CreateUser(&types.CreateUserPayload{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}