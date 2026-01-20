package infra

import (
	"context"
	"log"
	"shortvid-backend/app/shortvid-service/internal/conf"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func NewFirebaseApp(c *conf.Firebase) *auth.Client {
	opt := option.WithCredentialsJSON([]byte(c.CredentialsJson))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	log.Printf("Firebase admin sdk init success...")
	return client
}
