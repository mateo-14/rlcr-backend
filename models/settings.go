package models

type Settings struct {
	MaxSell         int             `json:"maxSell,omitempty"`
	MaxBuy          int             `json:"maxBuy,omitempty"`
	CreditSellValue float32         `json:"creditSellValue,omitempty"`
	CreditBuyValue  float32         `json:"creditBuyValue,omitempty"`
	BuyEnabled      bool            `json:"buyEnabled,omitempty"`
	SellEnabled     bool            `json:"sellEnabled,omitempty"`
	PaymentMethods  []PaymentMethod `json:"paymentMethods,omitempty"`
}

func (m *Settings) ToMap() (map[string]interface{}, error) {
	return toMap(m)
}
