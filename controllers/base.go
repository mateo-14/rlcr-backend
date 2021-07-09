package controllers

import (
	"fmt"

	"github.com/gorilla/mux"
)

type IController interface {
	Handle(prefix string, router *mux.Router)
}

type Controller struct {
	Prefix string
	router *mux.Router
}

func (c *Controller) Handle(prefix string, router *mux.Router) {
	prefix = fmt.Sprintf("/%s%s", prefix, c.Prefix)
	c.router = router.PathPrefix(prefix).Subrouter()
}
