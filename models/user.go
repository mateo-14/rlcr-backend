package models

import (
	"github.com/diamondburned/arikawa/discord"
)

type User struct {
	ID           discord.UserID `json:"-"`
	RefreshToken string         `json:"refreshToken,omitempty"`
}

func (m *User) ToMap() (map[string]interface{}, error) {
	return toMap(m)
}
