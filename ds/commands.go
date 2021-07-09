package ds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nico-14/rlcr-backend/models"
	"github.com/Nico-14/rlcr-backend/models/orderm"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

func (b *Bot) Pedidos(m *gateway.MessageCreateEvent) (string, error) {
	embed := &discord.Embed{
		Timestamp: discord.Timestamp(time.Now()),
		Color:     discord.Color(0x8B5CF6),
	}

	token, err := b.s.UsrSvc.GenerateOrdersToken(b.Context.Context(), m.Author.ID)
	if err != nil {
		return "", errors.New("")
	}

	orders, err := b.s.UsrSvc.GetOrders(b.Context.Context(), m.Author.ID, 5, "")

	if err == nil {
		for i := range orders {
			order := &orders[i]
			var value string
			if order.Mode == 1 {
				value = "Venta de"
			} else {
				value = "Compra de"
			}

			embed.Fields = append(embed.Fields, discord.EmbedField{
				Name:  fmt.Sprintf("Pedido %s", order.ID),
				Value: fmt.Sprintf("[%s %v créditos a ARS$ %v](%s/orders/%s?t=%s)", value, order.Credits, order.Price, os.Getenv("FRONTEND_URL"), order.ID, token),
			})
		}
	}

	if err != nil || len(orders) == 0 {
		embed.Title = "No has realizado ningún pedido :frowning:"
	} else {
		embed.Title = "Ver todos los pedidos"
		embed.URL = fmt.Sprintf("%s/orders?t=%s", os.Getenv("FRONTEND_URL"), token)
		embed.Description = "Últimos pedidos"
	}

	b.SendMessage(m.ChannelID, "", embed)

	return "", errors.New("")
}

func (b *Bot) waitResp(ctx context.Context, m *gateway.MessageCreateEvent) interface{} {
	r := b.WaitFor(ctx, func(i interface{}) bool {
		mg, ok := i.(*gateway.MessageCreateEvent)
		if !ok {
			return false
		}
		return mg.Author.ID == m.Author.ID
	})
	if r == nil {
		b.SendMessage(m.ChannelID, "Pedido cancelado ❌", nil)
	}
	return r
}

func (b *Bot) Creditos(m *gateway.MessageCreateEvent) (string, error) {
	settings, err := b.s.SettSvc.Get(b.Context.Context())
	if err != nil {
		return "", err
	}

	if _, err := b.SendMessage(m.ChannelID, ":thinking: Querés comprar o vender créditos? Responde con **comprar** o **vender**", nil); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r := b.waitResp(ctx, m)
	if r == nil {
		return "", errors.New("")
	}

	ev := r.(*gateway.MessageCreateEvent)
	if ev.Content != "comprar" && ev.Content != "vender" {
		return ":thinking: Mmm... no entendí tu respuesta", nil
	}

	mode := strings.ToLower(ev.Content)
	order := orderm.Order{}
	if mode == "comprar" {
		order.Mode = orderm.Buy
	} else {
		order.Mode = orderm.Sell
	}

	var max int
	if order.Mode == orderm.Buy {
		max = settings.MaxSell
	} else {
		max = settings.MaxBuy
	}

	if _, err := b.SendMessage(m.ChannelID, fmt.Sprintf(":thinking: Cuántos créditos vas a %s? (Mínimo **100**, máximo **%v**)", mode, max), nil); err != nil {
		return "", err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r = b.waitResp(ctx, m)
	if r == nil {
		return "", errors.New("")
	}

	ev = r.(*gateway.MessageCreateEvent)
	cr, err := strconv.ParseInt(ev.Content, 10, 32)
	if err != nil {
		return ":face_with_raised_eyebrow: Todavía no soy lo suficientemente inteligente como para reconocer ese \"número\"", nil
	}

	if cr < 100 {
		return fmt.Sprintf(":neutral_face: El mínimo de créditos para %s es de 100", mode), nil
	}

	if int(cr) > max {
		return fmt.Sprintf(":neutral_face: El máximo de créditos para %s es de %v", mode, max), nil
	}

	order.Credits = int(cr)
	order.Sanitize(settings)

	msg, err := b.SendMessage(m.ChannelID, fmt.Sprintf(":dollar: Vas a %s %v créditos a ARS$ %v. Reacciona al mensaje con ✅ para confirmar tu pedido", mode, order.Credits, order.Price), nil)
	if err != nil {
		return "", err
	}

	err = b.React(m.ChannelID, msg.ID, "✅")
	if err != nil {
		return "", err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r = b.WaitFor(ctx, func(i interface{}) bool {
		reac, ok := i.(*gateway.MessageReactionAddEvent)
		if !ok {
			return false
		}
		return reac.Emoji.APIString() == "✅" && reac.UserID == m.Author.ID
	})

	if r != nil {
		err := b.s.UsrSvc.AddOrder(b.Context.Context(), &models.User{ID: m.Author.ID}, &order)
		if err != nil {
			return fmt.Sprintf("🥴 Ocurrió un error al realizar el pedido. Error: %s", err.Error()), nil
		}

		t, err := b.s.UsrSvc.GenerateOrdersToken(b.Context.Context(), m.Author.ID)
		if err != nil {
			return fmt.Sprintf("🥴 Ocurrió un error al realizar el pedido. Error: %s", err.Error()), nil
		}

		return fmt.Sprintf("Pedido confirmado ✅\n%s/orders/%s?t=%s", os.Getenv("FRONTEND_URL"), order.ID, t), nil
	} else {
		b.SendMessage(m.ChannelID, "Pedido cancelado ❌", nil)
		b.Unreact(m.ChannelID, msg.ID, "✅")
	}

	return "", errors.New("")
}
