package typeDefines

import (
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type IAuth interface {
	authenticate() (string, error)
}

type oAuth2 struct {
	Grant_Type    string
	Token_URL     string
	CLIENT_ID     string
	CLIENT_SECRET string
	Header_Prefix string
	Token         string
}

func NewoAuth2(grant_type string) oAuth2 {
	var oauth2 oAuth2
	oauth2.Grant_Type = grant_type
	return oauth2
}

func (oAuth2 *oAuth2) Authenticate() (string, error) {
	oAuth2.collectDotEnv()

	baseURL, err := url.Parse(oAuth2.Token_URL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("grant_type", oAuth2.Grant_Type)
	queryParams.Add("client_id", oAuth2.CLIENT_ID)
	queryParams.Add("client_secret", oAuth2.CLIENT_SECRET)

	baseURL.RawQuery = queryParams.Encode()

	request, err := http.NewRequest("POST", baseURL.String(), strings.NewReader(""))
	if err != nil {
		log.Fatalf("An error ocurred while creating the request %v\n", err)
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("An error ocurred while making the request %v\n", err)
		return "", err
	}

	//puts response.Body in a []byte
	out, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("An error ocurred when reading the response %v\n", err)
		return "", err

	}
	token_str := string(out)
	return token_str, nil
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
