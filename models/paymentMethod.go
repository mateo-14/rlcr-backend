package models

type PaymentMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}
