package eventinfo

import (
	"context"
	"event-history/pkg/eventinfo/dto"
	"event-history/pkg/eventinfo/model"
	"event-history/pkg/repository"
	"fmt"
)

type Service interface {
	CreateKey(ctx context.Context, info *model.EventSnapshot) error
	GetAnswer(ctx context.Context, eventsQuery *dto.EventQuery) (*dto.EventResponse, error)
	DeleteKey(ctx context.Context, e *dto.EventQuery) error
	UpdateKey(ctx context.Context, m *model.EventSnapshot) error
	GetHistory(ctx context.Context, e *dto.EventQuery) ([]dto.EventHistoryResponse, error)
}

type EventService struct {
	repository repository.EventRepository
}

func (es *EventService) CreateKey(ctx context.Context, info *model.EventSnapshot) error {
	err := es.repository.CreateKey(ctx, info)
	if err != nil {
		return fmt.Errorf("Service.CreateKey failed. Error: %w", err)
	}

	return nil
}

func (es *EventService) UpdateKey(ctx context.Context, info *model.EventSnapshot) error {
	err := es.repository.UpdateKey(ctx, info)
	if err != nil {
		return fmt.Errorf("Service.UpdateKey failed. Error: %w", err)
	}

	return nil
}

func (es *EventService) GetAnswer(ctx context.Context, eventQuery *dto.EventQuery) (*dto.EventResponse, error) {
	eventInfo, err := es.repository.GetAnswer(ctx, eventQuery)
	if err != nil {
		return nil, fmt.Errorf("Service.GetAnswer: %+v", err)
	}

	return dto.NewEventResponse(eventInfo), nil
}

func (es *EventService) DeleteKey(ctx context.Context, eventQuery *dto.EventQuery) error {
	err := es.repository.DeleteKey(ctx, eventQuery)
	if err != nil {
		return fmt.Errorf("Service.DeleteKey: %+v", err)
	}
	return nil
}

func (es *EventService) GetHistory(ctx context.Context, eventQuery *dto.EventQuery) ([]dto.EventHistoryResponse, error) {
	history, err := es.repository.GetHistory(ctx, eventQuery)
	if err != nil {
		return nil, fmt.Errorf("Service.GetHistory: %+v", err)
	}

	return dto.NewEventHistoryResponse(history), nil
}

func NewEventService(repository repository.EventRepository) Service {
	return &EventService{
		repository: repository,
	}
}
