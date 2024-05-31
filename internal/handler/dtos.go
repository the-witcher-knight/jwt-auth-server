package handler

type generateTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type generateTokenResponse struct {
	AccessToken string `json:"access_token"`
}
