package domain

type Task struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Points      int    `gorm:"not null" json:"points"`
}

func (Task) TableName() string {
	return "tasks"
}

type UserTask struct {
	ID     int `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID int `gorm:"not null;index" json:"user_id"`
	TaskID int `gorm:"nott null;index" json:"task_id"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Task Task `gorm:"foreignKey:TaskID" json:"-"`
}

func (UserTask) TableName() string {
	return "user_tasks"
}
