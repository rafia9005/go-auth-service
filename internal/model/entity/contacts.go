package entity

type Contacts struct {
	ID     uint   `gorm:"primaryKey"`
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"`
	User   Users  `gorm:"foreignKey:UserID"`
}

