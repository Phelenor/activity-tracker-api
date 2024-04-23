package models

type LoginRequest struct {
	IdToken string `json:"idToken"`
	Nonce   string `json:"nonce"`
}

type UserTokenResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
