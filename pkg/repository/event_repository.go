package repository

import (
	"context"
	"event-history/pkg/eventinfo/dto"
	"event-history/pkg/eventinfo/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

//go:generate moq -out mock/EventRepository.go -pkg mock . EventRepository
type EventRepository interface {
	CreateKey(ctx context.Context, eventInfo *model.EventSnapshot) error
	GetAnswer(ctx context.Context, eventQuery *dto.EventQuery) (*model.EventSnapshot, error)
	DeleteKey(ctx context.Context, query *dto.EventQuery) error
	UpdateKey(ctx context.Context, info *model.EventSnapshot) error
	GetHistory(ctx context.Context, query *dto.EventQuery) ([]model.EventHistory, error)
}

type gormEventRepository struct {
	db *gorm.DB
}

func (gbr *gormEventRepository) CreateKey(ctx context.Context, eventInfo *model.EventSnapshot) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	tx := gbr.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	queryResult := gbr.db.WithContext(ctx).Create(&eventInfo)
	if queryResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("create key for: %s key for %s user failed. error %w", eventInfo.Key, eventInfo.UserId, queryResult.Error)
	}

	queryResult = tx.WithContext(ctx).Create(model.NewHistoryRecord(eventInfo, model.CreateAction))
	if queryResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to historize create event for %s/%s, error: %+v", eventInfo.Key, eventInfo.UserId, queryResult.Error)
	}

	tx.Commit()
	return nil
}

func (gbr *gormEventRepository) UpdateKey(ctx context.Context, eventInfo *model.EventSnapshot) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	tx := gbr.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.WithContext(ctx).Where("key = ? and user_id = ?", eventInfo.Key, eventInfo.UserId).Updates(eventInfo)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("update key for: %s key for %s user failed. error %+v", eventInfo.Key, eventInfo.UserId, result.Error)
	} else if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("update key for: %s key for %s not found", eventInfo.Key, eventInfo.UserId)
	}

	result = tx.WithContext(ctx).Create(model.NewHistoryRecord(eventInfo, model.UpdateAction))
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to historize update event for %s/%s, error: %+v", eventInfo.Key, eventInfo.UserId, result.Error)
	}

	tx.Commit()

	return nil
}

func (gbr *gormEventRepository) GetAnswer(ctx context.Context, eventQuery *dto.EventQuery) (*model.EventSnapshot, error) {
	var res model.EventSnapshot
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	db := gbr.db.WithContext(ctx).Where("key = ? and user_id = ?", eventQuery.Key, eventQuery.UserId).First(&res)
	if db.Error != nil {
		return nil, fmt.Errorf("get answer for: %s key for %s user failed: %w", eventQuery.Key, eventQuery.UserId, db.Error)
	}

	return &res, nil
}

func (gbr *gormEventRepository) DeleteKey(ctx context.Context, eventquery *dto.EventQuery) error {
	var res model.EventSnapshot
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	tx := gbr.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	execResult := gbr.db.WithContext(ctx).Unscoped().Where("key = ? and user_id = ?", eventquery.Key, eventquery.UserId).Delete(&res)
	if execResult.Error != nil {
		return fmt.Errorf("delete key for: %s key for %s user failed: %w", eventquery.Key, eventquery.UserId, execResult.Error)
	} else if execResult.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("record not found for %s key %s user", eventquery.Key, eventquery.UserId)
	}

	execResult = tx.WithContext(ctx).Create(model.NewHistoryRecord(
		&model.EventSnapshot{Key: eventquery.Key, UserId: eventquery.UserId},
		model.DeleteAction),
	)
	if execResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to historize delete event for %s/%s, error: %+v", eventquery.Key, eventquery.UserId, execResult.Error)
	}

	tx.Commit()

	return nil
}

func (gbr *gormEventRepository) GetHistory(ctx context.Context, eventQuery *dto.EventQuery) ([]model.EventHistory, error) {
	var res []model.EventHistory
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	db := gbr.db.WithContext(ctx).Where("key = ? and user_id = ?", eventQuery.Key, eventQuery.UserId).Order("created_at").Find(&res)
	if db.Error != nil {
		return nil, fmt.Errorf("failed to get history for %s/%s, error: %+v", eventQuery.UserId, eventQuery.Key, db.Error)
	}

	return res, nil
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &gormEventRepository{
		db: db,
	}
}
