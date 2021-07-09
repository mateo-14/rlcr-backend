package services

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Nico-14/rlcr-backend/db"
	"github.com/Nico-14/rlcr-backend/models"
	"github.com/Nico-14/rlcr-backend/models/orderm"
	"github.com/diamondburned/arikawa/discord"
	gonanoid "github.com/matoous/go-nanoid"
)

type IUsersService interface {
	AddOrder(context.Context, *models.User, *orderm.Order) error
	GetOrders(ctx context.Context, uid discord.UserID, limit int, startAfter string) (orders []orderm.Order, err error)
	GenerateOrdersToken(ctx context.Context, uid discord.UserID) (token string, err error)
	GetOrdersByToken(ctx context.Context, token string, limit int, startAfter string) (orders []orderm.Order, err error)
	getUserIdByToken(ctx context.Context, token string) (userID discord.UserID, err error)
	GetOrderByToken(ctx context.Context, token string, orderID string) (order *orderm.Order, userID discord.UserID, err error)
}

type UsersServices struct {
	db *db.Client
}

func NewUsersService(db *db.Client) *UsersServices {
	return &UsersServices{db}
}

func (s *UsersServices) AddOrder(ctx context.Context, user *models.User, order *orderm.Order) error {
	ref := s.db.Collection("users").Doc(user.ID.String())
	userm, err := user.ToMap()
	if err != nil {
		return err
	}
	if _, err := ref.Set(ctx, userm, firestore.MergeAll); err != nil {
		return err
	}
	ref.Collection("orders").Doc(order.ID).Set(ctx, order)
	return nil
}

func (s *UsersServices) GetOrders(ctx context.Context, uid discord.UserID, limit int, startAfter string) ([]orderm.Order, error) {
	if limit == 0 {
		limit = 10
	}

	ref := s.db.Collection("users").Doc(uid.String()).Collection("orders")
	q := ref.OrderBy("createdAt", firestore.Desc)

	if startAfter != "" {
		doc, _ := ref.Doc(startAfter).Get(ctx)
		q = q.StartAfter(doc)
	}

	docs, err := q.Limit(limit).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	orders := make([]orderm.Order, len(docs))
	for i := range docs {
		docs[i].DataTo(&orders[i])
		orders[i].ID = docs[i].Ref.ID
	}
	return orders, nil
}

func (s *UsersServices) GenerateOrdersToken(ctx context.Context, uid discord.UserID) (string, error) {
	exp := time.Now().Add(time.Minute * 30)
	nidl, err := strconv.ParseInt(os.Getenv("NANOID_ORDERS_TOKEN_LENGTH"), 10, 32)
	if err != nil {
		nidl = 18
	}

	token, err := gonanoid.Generate(os.Getenv("NANOID_ORDERS_TOKEN_ALPHABET"), int(nidl))
	if err != nil {
		return "", err
	}

	ref := s.db.Collection("orderstokens").Doc(token)
	if _, err = ref.Set(ctx, map[string]interface{}{"userID": uid.String()}); err != nil {
		return "", err
	}

	if err = s.db.AddExpireDoc(ctx, ref, exp); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UsersServices) getUserIdByToken(ctx context.Context, token string) (discord.UserID, error) {
	dsnap, err := s.db.Collection("orderstokens").Doc(token).Get(ctx)
	if err != nil {
		return 0, errors.New("expired or incorrect token")
	}

	tokend := dsnap.Data()
	uid, err := strconv.ParseInt(tokend["userID"].(string), 10, 64)
	return discord.UserID(uid), err
}

func (s *UsersServices) GetOrdersByToken(ctx context.Context, token string, limit int, startAfter string) ([]orderm.Order, error) {
	uid, err := s.getUserIdByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.GetOrders(ctx, discord.UserID(uid), limit, startAfter)
}

func (s *UsersServices) GetOrderByToken(ctx context.Context, token string, orderID string) (*orderm.Order, discord.UserID, error) {
	uid, err := s.getUserIdByToken(ctx, token)
	if err != nil {
		return nil, 0, err
	}

	dsnap, err := s.db.Collection("users").Doc(uid.String()).Collection("orders").Doc(orderID).Get(ctx)
	if err != nil {
		return nil, 0, err
	}

	var order orderm.Order
	err = dsnap.DataTo(&order)
	if err == nil {
		order.ID = dsnap.Ref.ID

	}
	return &order, uid, err
}
