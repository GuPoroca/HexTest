package typeDefines

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
)

type IAuth interface {
	Authenticate() (string, error)
}

type oAuth2 struct {
	Grant_Type    string
	Token_URL     string
	CLIENT_ID     string
	CLIENT_SECRET string
	Scopes        []string
	Token         string
}

func NewoAuth2(grant_type string) *oAuth2 {
	var oauth2 oAuth2
	oauth2.Grant_Type = grant_type
	return &oauth2
}

func (oAuth2 oAuth2) Authenticate() (string, error) {
	ctx := context.Background()
	oAuth2.collectDotEnv()
	auth := clientcredentials.Config{ClientID: oAuth2.CLIENT_ID, ClientSecret: oAuth2.CLIENT_SECRET, TokenURL: oAuth2.Token_URL, Scopes: oAuth2.Scopes}
	token_str, err := auth.Token(ctx)
	if err != nil {
		log.Fatalf("error handling auth, %v", err)
	}
	return token_str.AccessToken, err
}

func (oAuth2 *oAuth2) collectDotEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	oAuth2.CLIENT_ID = os.Getenv("CLIENT_ID")
	oAuth2.CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	oAuth2.Token_URL = os.Getenv("Token_URL")

}
