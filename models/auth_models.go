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

type UpdateUserDataRequest struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BirthTimestamp int64  `json:"birthTimestamp"`
}
