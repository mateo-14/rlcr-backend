package orderm

import (
	"os"
	"strconv"
	"time"

	"github.com/Nico-14/rocket-credits-backend/models"
	"github.com/Nico-14/rocket-credits-backend/util"
	gonanoid "github.com/matoous/go-nanoid"
)

type OrderStatus int

const (
	WaitingConfirmation OrderStatus = iota
	WaitingPayment
	Completed
)

type OrderMode int

const (
	Buy OrderMode = iota
	Sell
)

type Order struct {
	Mode            OrderMode   `json:"mode" firestore:"mode"`
	PaymentMethodId int         `json:"paymentMethodId" firestore:"paymentMethodId"`
	Credits         int         `json:"credits,omitempty" firestore:"credits,omitempty"`
	Price           float32     `json:"price,omitempty" firestore:"price,omitempty"`
	Code            string      `json:"code,omitempty" firestore:"-"`
	CreatedAt       time.Time   `json:"createdAt,omitempty" firestore:"createdAt"`
	Status          OrderStatus `json:"status" firestore:"status"`
	ID              string      `json:"id,omitempty" firestore:"-"`
	DNI             string      `json:"dni,omitempty" firestore:"dni,omitempty"`
	Account         string      `json:"account,omitempty" firestore:"account,omitempty"`
	PaymentAccount  string      `json:"paymentAccount,omitempty" firestore:"paymentAccount,omitempty"`
	Cvu             string      `json:"cvu,omitempty" firestore:"cvu,omitempty"`
	DsUsername      string      `json:"username,omitempty" firestore:"-"`
}

func (m *Order) Sanitize(settings *models.Settings) {
	nidl, err := strconv.ParseInt(os.Getenv("NANOID_ORDER_LENGTH"), 10, 32)
	if err != nil {
		nidl = 10
	}

	m.CreatedAt = time.Now()
	m.ID, _ = gonanoid.Generate(os.Getenv("NANOID_ORDER_ALPHABET"), int(nidl))
	pmExists := false
	for i := range settings.PaymentMethods {
		if settings.PaymentMethods[i].ID == m.PaymentMethodId {
			pmExists = true
			break
		}
	}

	if !pmExists {
		m.PaymentMethodId = 0
	}

	if m.Mode < 0 || m.Mode > 1 {
		m.Mode = 0
	}

	if m.Mode == 0 {
		m.Credits = util.Min(util.RoundCredits(m.Credits), settings.MaxSell)
		m.Price = float32(m.Credits) * settings.CreditSellValue
	} else {
		m.Credits = util.Min(util.RoundCredits(m.Credits), settings.MaxBuy)
		m.Price = float32(m.Credits) * settings.CreditBuyValue
	}
}
