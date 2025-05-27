package auth

import (
	"context"
	"os"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

type Auth struct {
	Pubkey *keyfunc.JWKS
	Config *oauth2.Config
}

var azureADConfig = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	RedirectURL:  "http://localhost:3000/callback",
	Scopes:       []string{"openid", "profile", "email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/oauth2/v2.0/token",
	},
}

func NewAuth() (*Auth, error) {
	pubkey, err := setupAzureJWTVerification()
	if err != nil {
		return nil, err
	}
	auth := &Auth{
		Pubkey: pubkey,
		Config: azureADConfig,
	}
	return auth, nil
}

func setupAzureJWTVerification() (*keyfunc.JWKS, error) {
	jwksURL := "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/discovery/v2.0/keys"
	return keyfunc.Get(jwksURL, keyfunc.Options{})
}

func (*Auth) VerifyToken(tokenStr string, jwks *keyfunc.JWKS) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, jwks.Keyfunc)
}

func (*Auth) RefreshAccessToken(token *oauth2.Token, config *oauth2.Config) (*oauth2.Token, error) {
	ctx := context.Background()
	ts := config.TokenSource(ctx, token)
	newToken, err := ts.Token() // Utilise le refresh_token automatiquement
	if err != nil {
		return nil, err
	}
	return newToken, nil
}
