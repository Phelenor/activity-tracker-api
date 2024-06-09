package gym

type GymEquipment struct {
	Id           string `json:"id" gorm:"primaryKey"`
	OwnerId      string `json:"ownerId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl"`
	VideoUrl     string `json:"videoUrl,omitempty"`
	ActivityType string `json:"activityType"`
}

type CreateEquipmentRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageUrl     string `json:"imageUrl"`
	VideoUrl     string `json:"videoUrl,omitempty"`
	ActivityType string `json:"activityType"`
}
