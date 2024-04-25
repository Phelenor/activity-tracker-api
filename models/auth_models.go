package models

type LoginRequest struct {
	IdToken string `json:"idToken"`
	Nonce   string `json:"nonce"`
}

type TokenRefreshRequest struct {
	Id string `json:"id"`
}

type UserTokenResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type UserChangeNameRequest struct {
	Name string `json:"name"`
}
