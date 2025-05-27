package auth_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Tomy2e/cluster-api-provider-scaleway/internal/auth"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
)

var testtoken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IkNOdjBPSTNSd3FsSEZFVm5hb01Bc2hDSDJYRSIsImtpZCI6IkNOdjBPSTNSd3FsSEZFVm5hb01Bc2hDSDJYRSJ9.eyJhdWQiOiJodHRwczovL21hbmFnZW1lbnQuYXp1cmUuY29tLyIsImlzcyI6Imh0dHBzOi8vc3RzLndpbmRvd3MubmV0L2FmMjJjNmVjLWJiNTktNDZjMS1iMTRlLWIwMzM1ZWU5NDdhMi8iLCJpYXQiOjE3NDgzMzg5OTYsIm5iZiI6MTc0ODMzODk5NiwiZXhwIjoxNzQ4MzQyODk2LCJhaW8iOiJrMlJnWVBpc3pQUitLMzlmOHZIOVdZa2k0amJsQUE9PSIsImFwcGlkIjoiYTQ4Nzk0NGMtZTZlNS00NzJiLWEzOTQtMjc2ZmZmZGE5NDBhIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvYWYyMmM2ZWMtYmI1OS00NmMxLWIxNGUtYjAzMzVlZTk0N2EyLyIsImlkdHlwIjoiYXBwIiwib2lkIjoiZGMxZDdjMWYtNjJiOC00M2UzLTk2NmUtY2E0M2E4M2IwYTQ3IiwicmgiOiIxLkFVc0E3TVlpcjFtN3dVYXhUckF6WHVsSG9rWklmM2tBdXRkUHVrUGF3ZmoyTUJNWkFRQkxBQS4iLCJzdWIiOiJkYzFkN2MxZi02MmI4LTQzZTMtOTY2ZS1jYTQzYTgzYjBhNDciLCJ0aWQiOiJhZjIyYzZlYy1iYjU5LTQ2YzEtYjE0ZS1iMDMzNWVlOTQ3YTIiLCJ1dGkiOiI3M2pjS09QYWprMi1hYTVkNzRaSkFBIiwidmVyIjoiMS4wIiwieG1zX2Z0ZCI6IklXRkRaTFM5bjY5NU80VjlWZDJmendtTVd6dllYMS1XYmwtRkp2c1hOb2tCWlhWeWIzQmxkMlZ6ZEMxa2MyMXoiLCJ4bXNfaWRyZWwiOiI3IDE4IiwieG1zX3JkIjoiMC40MkxsWUJKaTlCWVM0V0FYRXVBTjJMN3ZKdmQtNTIzWDV6MDQ2YkV6RVNqS0tTVGctbEt5TE9aamx2TzhSMTluVGRfeThoTlFsRU5Jd00yanZwSmgxM1hfWGM2M0R2NzRaaUFNQUEiLCJ4bXNfdGNkdCI6MTY2ODAwNTUwMn0.WJiAcHs9kiyhY5ZOc0QIKDvS4jN_mTsaBqrLCh_1eyOjpULEpYcbecngK9o_hgj6DY529RUHhuHaMj7VKHlMMuPt6CgikNq93Nkj415hULoKoB0449suKLg54pPhvFsQOAX4oDRubyz0eEZjss-H3-Ml9VdZUvyW_WxxwLhWPH7cVn8h75w03lkeVaJWlQa_R_DRIuDjdkPkst_DEC5neNdxPJzuuCAes1l71d45P-fwdIolCO78OltlstVJHivLiYw4-wZNQadrgN2P2k4Mki1XYBZ3pxcJYaBj-RLtxxteQBQFtrWvmOPebwIUr5_1FokQjZjKC0AiOebzMbDs_A"

func TestAuth(t *testing.T) {
	auth, err := auth.NewAuth()
	require.NoError(t, err)

	conf := &clientcredentials.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/oauth2/v2.0/token",
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	token, err := conf.Token(context.Background())
	if err != nil {
		panic(err)
	}

	verifiedToken, err := auth.VerifyToken(testtoken, auth.Pubkey)
	require.NoError(t, err)
	fmt.Printf("Token : %v\n", token)
	require.True(t, verifiedToken.Valid)
}
