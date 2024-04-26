package user

import (
	"fmt"
	"net/http"

	"github.com/Qrulaziz/jelasin/backend/services/auth"
	"github.com/Qrulaziz/jelasin/backend/types"
	"github.com/Qrulaziz/jelasin/backend/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// check if the user exist
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User with email %s already exists", payload.Email))
		return
	}

	// Hashing password
	hashedPassword, err := auth.HashPassword(payload.Password)

	// if it doesnt we create the new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
