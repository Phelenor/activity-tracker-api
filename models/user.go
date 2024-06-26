package models

type User struct {
	Id             string `json:"id" gorm:"primaryKey"`
	DisplayName    string `json:"display_name"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ImageUrl       string `json:"imageUrl,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	Height         int    `json:"height,omitempty"`
	BirthTimestamp int64  `json:"birthTimestamp,omitempty"`
}
