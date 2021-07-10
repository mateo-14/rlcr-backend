package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nico-14/rlcr-backend/db"
	"github.com/Nico-14/rlcr-backend/middlewares"
	"github.com/gorilla/mux"
)

type AuthBody struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthController struct {
	*Controller
	dbClient *db.Client
}

func NewAuthController(prefix string, dbClient *db.Client) *AuthController {
	c := &AuthController{dbClient: dbClient}
	c.Controller = &Controller{Prefix: prefix}
	return c
}

func (c *AuthController) Handle(prefix string, router *mux.Router) {
	c.Controller.Handle(prefix, router)
	c.router.HandleFunc("/", c.auth).Methods(http.MethodPost, http.MethodOptions)
	c.router.HandleFunc("/", middlewares.VerifyToken(c.authWithToken)).Methods(http.MethodGet)
}

func (c *AuthController) auth(w http.ResponseWriter, r *http.Request) {
	var authBody AuthBody
	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	if authBody.Password == nil || authBody.Username == nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	dsnap, err := c.dbClient.Collection("settings").Doc("user").Get(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	m := dsnap.Data()
	if hash, found := m["password"]; found {
		if verifyPassword(*authBody.Password, hash.(string)) {
			token, err := generateToken(time.Hour)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(&AuthResponse{Token: token})
				return
			}
		}
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (c *AuthController) authWithToken(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken(time.Hour)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&AuthResponse{Token: token})
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
