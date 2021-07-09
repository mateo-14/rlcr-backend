package db

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
	quit chan struct{}
}

func New() *Client {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	log.Println("Firebase connected!")

	c := &Client{client, make(chan struct{})}

	go c.check()
	return c
}

type ExpirationDoc struct {
	Doc        *firestore.DocumentRef `firestore:"doc"`
	ExpireDate time.Time              `firestore:"expireDate"`
}

func (c *Client) check() {
checkLoop:
	for {
		select {
		case <-c.quit:
			break checkLoop
		default:
			if docs, err := c.Client.Collection("expiration").Documents(context.Background()).GetAll(); err == nil {
				batch := c.Batch()

				var expDoc *ExpirationDoc
				for i := range docs {
					docs[i].DataTo(&expDoc)
					if expDoc != nil && expDoc.ExpireDate.Unix() <= time.Now().Unix() {
						batch.Delete(docs[i].Ref)
						batch.Delete(expDoc.Doc)
					}
					expDoc = nil
				}

				batch.Commit(context.Background())
			}
			time.Sleep(time.Minute)
		}
	}
}

func (c *Client) Close() {
	c.quit <- struct{}{}
	c.Client.Close()
}

func (c *Client) AddExpireDoc(ctx context.Context, docRef *firestore.DocumentRef, expireDate time.Time) error {
	expDoc := &ExpirationDoc{Doc: docRef, ExpireDate: expireDate}
	_, _, err := c.Client.Collection("expiration").Add(ctx, expDoc)
	return err
}
