package ds

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nico-14/rlcr-backend/services"
	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
)

type Bot struct {
	*bot.Context
	s *services.Services
}

var Client *Bot

func Connect(services *services.Services) {
	Client = &Bot{s: services}
	bot.Start(os.Getenv("CLIENT_TOKEN"), Client, func(c *bot.Context) error {
		c.HasPrefix = bot.NewPrefix("!")
		return nil
	})

	bot.UnknownCommandString = func(err *bot.ErrUnknownCommand) string {
		return ""
	}

	// go Client.updateActivity()
}

func (b *Bot) updateActivity() {
	for {
		// b.PresenceSet(nil, discord.Presence{User: b.Ready.User, Game: &discord.Activity{}})
		time.Sleep(time.Minute)
	}
}

// getOAuthToken returns token and refreshToken
func (b *Bot) getOAuthToken(codetoken string, isRefreshToken bool, redirectUrl string) (string, string, error) {
	form := url.Values{}

	form.Set("client_id", os.Getenv("CLIENT_ID"))
	form.Set("client_secret", os.Getenv("CLIENT_SECRET"))

	if isRefreshToken {
		form.Set("refresh_token", codetoken)
		form.Set("grant_type", "refresh_token")
	} else {
		form.Set("code", codetoken)
		form.Set("grant_type", "authorization_code")
		form.Set("redirect_uri", redirectUrl)
	}

	if resp, err := http.Post(fmt.Sprintf("%v/oauth2/token", os.Getenv("DS_API_ENDPOINT")), "application/x-www-form-urlencoded", strings.NewReader(form.Encode())); err != nil {
		return "", "", err
	} else {
		var bodyr map[string]interface{}
		defer resp.Body.Close()

		bodyb, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}

		if err := json.Unmarshal(bodyb, &bodyr); err != nil {
			return "", "", err
		}
		if t, ok := bodyr["access_token"]; ok {
			return fmt.Sprintf("%v", t), fmt.Sprintf("%v", bodyr["refresh_token"]), nil
		} else {
			return "", "", errors.New("no token found")
		}
	}

}

func (b *Bot) GetUserIDByToken(token string) (discord.UserID, error) {
	r, err := b.NewRequest(b.Client.Context(), "GET", fmt.Sprintf("%v/users/@me", os.Getenv("DS_API_ENDPOINT")))
	if err != nil {
		return 0, err
	}

	r.AddHeader(http.Header{"Authorization": {fmt.Sprintf("Bearer %v", token)}, "Content-Type": {"application/json"}})

	if resp, err := b.Do(r); err != nil {
		return 0, err
	} else {
		var bodyr map[string]interface{}

		body := resp.GetBody()
		defer body.Close()

		if err := json.NewDecoder(body).Decode(&bodyr); err != nil {
			return 0, err
		}
		if id, ok := bodyr["id"]; ok {
			id, err := strconv.ParseInt(id.(string), 10, 64)
			return discord.UserID(id), err
		}
		return 0, err
	}
}

// GetUserIdByCode returns UserID, token, refreshToken
func (b *Bot) GetUserIdByCode(code, redirectUrl string) (userID discord.UserID, token string, refreshToken string, _ error) {
	t, rt, err := b.getOAuthToken(code, false, redirectUrl)
	if err != nil {
		return 0, "", "", err
	}
	id, err := b.GetUserIDByToken(t)
	return id, t, rt, err
}

func (b *Bot) GetUserIdByRefreshToken(refreshToken, redirectUrl string) (discord.UserID, string, error) {
	t, _, err := b.getOAuthToken(refreshToken, true, redirectUrl)
	if err != nil {
		return 0, "", err
	}
	id, err := b.GetUserIDByToken(t)
	return id, t, err
}
