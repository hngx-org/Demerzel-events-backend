package services

import (
	"context"
	"demerzel-events/internal/oauth"
	"demerzel-events/pkg/types"
	"encoding/json"
	"log"
)

type OAuthCallbackParams struct {
	ctx   context.Context
	code  string
	state string
}

func NewOAuthCallbackParams(ctx context.Context, code, state string) OAuthCallbackParams {
	return OAuthCallbackParams{
		ctx: ctx, code: code, state: state,
	}
}

func OAuthCallback(params OAuthCallbackParams) (*types.UserInfo, error) {
	oauthConfig := oauth.OauthConfig()
	token, err := oauthConfig.Exchange(params.ctx, params.code)
	if err != nil {
		return nil, err
	}

	client := oauthConfig.Client(params.ctx, token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer userInfoResp.Body.Close()

	var userInfo types.UserInfo
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		log.Fatal("Error decoding user info:", err)
	}

	return &userInfo, nil
}
