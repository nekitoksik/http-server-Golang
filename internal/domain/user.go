package domain

type User struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"column:password_hash;not null" json:"-"`
	Balance      int    `gorm:"default:0" json:"balance"`
	RefreshToken string `gorm:"column:refresh_token" json:"-"`
	ReferrerID   *int   `gorm:"index" json:"referrer_id,omitempty"`

	Referrer       *User      `gorm:"foreignKey:ReferrerID" json:"-"`
	CompletedTasks []UserTask `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}
