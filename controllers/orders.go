package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Nico-14/rlcr-backend/ds"
	"github.com/Nico-14/rlcr-backend/models"
	"github.com/Nico-14/rlcr-backend/models/orderm"
	"github.com/Nico-14/rlcr-backend/services"
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersController struct {
	*Controller
	s *services.Services
}

func NewOrdersController(prefix string, services *services.Services) *OrdersController {
	c := &OrdersController{s: services}
	c.Controller = &Controller{Prefix: prefix}
	return c
}

func (c *OrdersController) Handle(prefix string, router *mux.Router) {
	c.Controller.Handle(prefix, router)
	c.router.HandleFunc("/", c.addOrder).Methods(http.MethodPost, http.MethodOptions)
	c.router.HandleFunc("/", c.getOrders).Methods(http.MethodGet)
	c.router.HandleFunc("/auth", c.generateToken).Methods(http.MethodGet)
	c.router.HandleFunc("/{oid}", c.getOrder).Methods(http.MethodGet)
}

func (c *OrdersController) addOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var order orderm.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	if order.Credits == 0 || order.Code == "" {
		http.Error(w, "Incomplete body", http.StatusBadRequest)
		return
	}

	settings, err := c.s.SettSvc.Get(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uid, token, rtoken, err := ds.Client.GetUserIdByCode(order.Code, fmt.Sprintf("%s/redirect", os.Getenv("FRONTEND_URL")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	svid, err := strconv.ParseUint(os.Getenv("SV_ID"), 10, 64)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := ds.Client.AddMember(discord.GuildID(svid), uid, api.AddMemberData{Token: token}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dmc, err := ds.Client.CreatePrivateChannel(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order.Sanitize(settings)
	if err := c.s.UsrSvc.AddOrder(r.Context(), &models.User{ID: uid, RefreshToken: rtoken}, &order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err = c.s.UsrSvc.GenerateOrdersToken(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var value string
	if order.Mode == 1 {
		value = "venta"
	} else {
		value = "compra"
	}

	url := fmt.Sprintf("%s/orders/%s?t=%s", os.Getenv("FRONTEND_URL"), order.ID, token)
	embed := &discord.Embed{
		Timestamp:   discord.Timestamp(time.Now()),
		Color:       discord.Color(0x8B5CF6),
		Title:       "Pedido realizado",
		URL:         url,
		Description: fmt.Sprintf("**__Has realizado un pedido de %s de %v créditos a ARS$ %v__**\n\n• El pedido debe ser confirmado por un moderador. Una vez confirmado nos contactaremos por DM para realizar la transacción.\n• Si tenés algún problema o necesitas ayuda usa el comando **!ayuda** o contacta con un moderador en nuestro [canal de discord](%s).\n• Usa el comando **!pedidos** para ver la lista con los últimos pedidos.\n\n[Abrir pedido](%s)", value, order.Credits, order.Price, os.Getenv("DS_CHANNEL"), url),
	}

	ds.Client.SendMessage(dmc.ID, "", embed)
	u, _ := ds.Client.User(uid)
	json.NewEncoder(w).Encode(map[string]string{"username": fmt.Sprintf("%s#%s", u.Username, u.Discriminator)})
}

func (c *OrdersController) getOrders(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	sa := r.URL.Query().Get("sa")

	if orders, err := c.s.UsrSvc.GetOrdersByToken(r.Context(), t, 10, sa); err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	} else {
		json.NewEncoder(w).Encode(orders)
	}
}

func (c *OrdersController) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := r.URL.Query().Get("t")
	oid := vars["oid"]

	if order, uid, err := c.s.UsrSvc.GetOrderByToken(r.Context(), t, oid); err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	} else {
		u, _ := ds.Client.User(uid)
		order.DsUsername = fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
		json.NewEncoder(w).Encode(order)
	}
}

func (c *OrdersController) generateToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	uid, _, _, err := ds.Client.GetUserIdByCode(code, "http://localhost:3000/orders")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if token, err := c.s.UsrSvc.GenerateOrdersToken(r.Context(), uid); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
