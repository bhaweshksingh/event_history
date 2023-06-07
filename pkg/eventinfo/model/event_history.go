package model

import (
	"time"
)

const (
	CreateAction = "create"
	UpdateAction = "update"
	DeleteAction = "delete"
)

type EventHistory struct {
	Key       string    `gorm:"column:key;" json:"key"`
	Value     string    `gorm:"column:value;" json:"value"`
	UserId    string    `gorm:"column:user_id" json:"user_id"`
	Action    string    `gorm:"column:action" json:"action"`
	CreatedAt time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
}

func (EventHistory) TableName() string {
	return "event_history"
}

func NewHistoryRecord(info *EventSnapshot, action string) *EventHistory {
	return &EventHistory{Key: info.Key, Value: info.Value, UserId: info.UserId, Action: action}
}
