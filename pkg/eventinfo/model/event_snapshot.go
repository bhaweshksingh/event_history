package model

type EventSnapshot struct {
	Key     string    `gorm:"column:key;" json:"key"`
	Value     string    `gorm:"column:value;" json:"value"`
	UserId    string    `gorm:"column:user_id" json:"user_id"`
}

func (EventSnapshot) TableName() string {
	return "event_snapshot"
}
