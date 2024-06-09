package gym

type Equipment struct {
	Id           string `json:"id" gorm:"primaryKey"`
	OwnerId      string `json:"ownerId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl,omitempty"`
	VideoUrl     string `json:"VideoUrl,omitempty"`
	ActivityType string `json:"activityType"`
}

type CreateEquipmentRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl,omitempty"`
	VideoUrl     string `json:"VideoUrl,omitempty"`
	ActivityType string `json:"activityType"`
}
