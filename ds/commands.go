package ds

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

func (b *Bot) Pedidos(m *gateway.MessageCreateEvent) (string, error) {
	dmc, err := b.CreatePrivateChannel(m.Author.ID)
	if err != nil {
		return "", err
	}

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

	_, err = b.SendMessage(dmc.ID, "", embed)

	if err == nil && m.ChannelID != dmc.ID {
		return "Esta información es privada, hablemos por DM... 🤫", nil
	}
	return "", errors.New("")
}

func (b *Bot) Ayuda(m *gateway.MessageCreateEvent) (string, error) {
	embed := &discord.Embed{
		Timestamp:   discord.Timestamp(time.Now()),
		Color:       discord.Color(0x8B5CF6),
		Title:       "Ayuda",
		Description: fmt.Sprintf("**__Lista de comandos__**\n• **!pedidos** **-** Muestra la lista con los últimos pedidos\n• **!info** **-** Muestra la información actual sobre el precio de los créditos y transacciones totales\n\n*Contacta con un moderador en nuestro [canal de discord](%s) si tenés algún problema o consulta.", os.Getenv("DS_CHANNEL_INVITE_URL")),
	}

	b.SendMessage(m.ChannelID, "", embed)
	return "", errors.New("")
}

func (b *Bot) Info(m *gateway.MessageCreateEvent) (string, error) {
	sts, err := b.s.SettSvc.Get(b.Context.Context())
	embed := &discord.Embed{
		Timestamp: discord.Timestamp(time.Now()),
		Color:     discord.Color(0x8B5CF6),
		Title:     "Información de la tienda",
	}

	if err != nil {
		embed.Description = "Hubo un problema al cargar la información"
	} else {
		var sellTxt, buyText string

		if sts.SellEnabled {
			sellTxt = fmt.Sprintf("100cr x $%v", sts.CreditSellValue*100)
		} else {
			sellTxt = "Venta deshabilitada"
		}

		if sts.BuyEnabled {
			buyText = fmt.Sprintf("100cr x $%v", sts.CreditBuyValue*100)
		} else {
			buyText = "Compra deshabilitada"
		}

		embed.Fields = []discord.EmbedField{
			{
				Name:  "Precio de venta",
				Value: sellTxt,
			},
			{
				Name:  "Precio de compra",
				Value: buyText,
			},
			{
				Name:  "Ventas realizadas",
				Value: "0",
			},
			{
				Name:  "Compras realizadas",
				Value: "0",
			},
		}
	}

	b.SendMessage(m.ChannelID, "", embed)
	return "", errors.New("")
}
