package models

type LoginRequest struct {
	IdToken string `json:"idToken"`
	Nonce   string `json:"nonce"`
}
