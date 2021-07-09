package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nico-14/rocket-credits-backend/middlewares"
	"github.com/Nico-14/rocket-credits-backend/models"
	"github.com/Nico-14/rocket-credits-backend/services"
	"github.com/gorilla/mux"
)

type SettingsController struct {
	*Controller
	s *services.Services
}

func NewSettingsController(prefix string, services *services.Services) *SettingsController {
	c := &SettingsController{s: services}
	c.Controller = &Controller{Prefix: prefix}
	return c
}

func (c *SettingsController) Handle(prefix string, router *mux.Router) {
	c.Controller.Handle(prefix, router)
	c.router.HandleFunc("/", c.getConfig).Methods(http.MethodGet)
	c.router.HandleFunc("/", middlewares.VerifyToken(c.updateConfig)).Methods(http.MethodPut)
}

func (c *SettingsController) getConfig(w http.ResponseWriter, r *http.Request) {
	if settings, err := c.s.SettSvc.Get(r.Context()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

func (c *SettingsController) updateConfig(w http.ResponseWriter, r *http.Request) {
	var settings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	if err := c.s.SettSvc.Update(r.Context(), &settings); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Settings updated")
}
