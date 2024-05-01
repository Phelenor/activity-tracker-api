package models

type LoginRequest struct {
	IdToken string `json:"idToken"`
	Nonce   string `json:"nonce"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserTokenResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserChangeNameRequest struct {
	Name string `json:"name"`
}
