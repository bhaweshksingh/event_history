package eventinfo

import (
	"context"
	"errors"
	"event-history/pkg/eventinfo/dto"
	"event-history/pkg/eventinfo/model"
	"event-history/pkg/repository/mock"
	"github.com/smartystreets/assertions"
	"testing"
)

const userId = "test_user"

func TestGormEventRepository_GetAnswer(t *testing.T) {
	ctx := context.Background()
	eventSnapshot := model.EventSnapshot{"name", "john", userId}
	repositoryMock := mock.EventRepositoryMock{
		GetAnswerFunc: func(ctx context.Context, eventQuery *dto.EventQuery) (*model.EventSnapshot, error) {
			return &eventSnapshot, nil
		}}

	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	actualSnapshot, err := service.GetAnswer(ctx, &event)

	assertions.So(err, assertions.ShouldBeNil)
	assertions.ShouldEqual(actualSnapshot, eventSnapshot)
}

func TestGormEventRepository_GetAnswer_fails(t *testing.T) {
	ctx := context.Background()
	mockError := errors.New("failed to get")
	repositoryMock := mock.EventRepositoryMock{
		GetAnswerFunc: func(ctx context.Context, eventQuery *dto.EventQuery) (*model.EventSnapshot, error) {
			return nil, mockError
		}}

	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	actualSnapshot, err := service.GetAnswer(ctx, &event)

	assertions.So(actualSnapshot, assertions.ShouldBeNil)
	assertions.ShouldContain(err.Error(), mockError.Error())
}

func TestGormEventRepository_GetHistory(t *testing.T) {
	ctx := context.Background()
	historyRecords := []model.EventHistory{
		{Key: "Name", Value: "John", Action: "create"},
		{Key: "Name", Value: "sam", Action: "update"},
	}
	repositoryMock := mock.EventRepositoryMock{
		GetHistoryFunc: func(ctx context.Context, query *dto.EventQuery) ([]model.EventHistory, error) {
			return historyRecords, nil
		}}

	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	historyResponse, err := service.GetHistory(ctx, &event)

	assertions.So(err, assertions.ShouldBeNil)
	assertions.ShouldEqual(len(historyResponse), 2)
	assertions.ShouldEqual(historyResponse[0], dto.EventHistoryResponse{dto.Data{"Name", "John"}, "create"})
	assertions.ShouldEqual(historyResponse[1], dto.EventHistoryResponse{dto.Data{"Name", "sam"}, "create"})
}

func TestGormEventRepository_GetHistory_fails(t *testing.T) {
	ctx := context.Background()
	mockError := errors.New("failed to get")
	repositoryMock := mock.EventRepositoryMock{
		GetHistoryFunc: func(ctx context.Context, query *dto.EventQuery) ([]model.EventHistory, error) {
			return nil, mockError
		}}

	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	historyResponse, err := service.GetHistory(ctx, &event)

	assertions.ShouldEqual(len(historyResponse), 0)
	assertions.So(err, assertions.ShouldBeNil)
	assertions.ShouldContain(err.Error(), mockError.Error())
}

func TestGormEventRepository_DeleteKey(t *testing.T) {
	ctx := context.Background()
	repositoryMock := mock.EventRepositoryMock{
		DeleteKeyFunc: func(ctx context.Context, query *dto.EventQuery) error {
			return nil
		},
	}
	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	err := service.DeleteKey(ctx, &event)

	assertions.So(err, assertions.ShouldBeNil)
}

func TestGormEventRepository_DeleteKey_fails(t *testing.T) {
	ctx := context.Background()
	mockError := errors.New("failed to delete")
	repositoryMock := mock.EventRepositoryMock{
		DeleteKeyFunc: func(ctx context.Context, query *dto.EventQuery) error {
			return mockError
		},
	}
	service := NewEventService(&repositoryMock)
	event := dto.EventQuery{"name", userId}

	err := service.DeleteKey(ctx, &event)

	assertions.ShouldContain(err.Error(), mockError.Error())
}

func TestGormEventRepository_UpdateKey(t *testing.T) {
	ctx := context.Background()
	repositoryMock := mock.EventRepositoryMock{
		UpdateKeyFunc: func(ctx context.Context, info *model.EventSnapshot) error {
			return nil
		},
	}

	service := NewEventService(&repositoryMock)
	updateEvent := model.EventSnapshot{"name", "john", userId}

	err := service.UpdateKey(ctx, &updateEvent)

	assertions.So(err, assertions.ShouldBeNil)
}

func TestGormEventRepository_UpdateKey_fails(t *testing.T) {
	ctx := context.Background()
	updateEvent := model.EventSnapshot{"name", "john", userId}
	mockError := errors.New("failed to update")
	repositoryMock := mock.EventRepositoryMock{
		UpdateKeyFunc: func(ctx context.Context, event *model.EventSnapshot) error {
			if *event == updateEvent {
				return mockError
			}
			return nil
		},
	}

	service := NewEventService(&repositoryMock)

	err := service.UpdateKey(ctx, &updateEvent)

	assertions.ShouldContain(err.Error(), mockError.Error())
}

func TestGormEventRepository_CreateKey(t *testing.T) {
	ctx := context.Background()
	repositoryMock := mock.EventRepositoryMock{
		CreateKeyFunc: func(ctx context.Context, info *model.EventSnapshot) error {
			return nil
		},
	}

	service := NewEventService(&repositoryMock)
	updateEvent := model.EventSnapshot{"name", "john", userId}

	err := service.CreateKey(ctx, &updateEvent)

	assertions.So(err, assertions.ShouldBeNil)
}

func TestGormEventRepository_CreateKey_fails(t *testing.T) {
	ctx := context.Background()
	createEvent := model.EventSnapshot{"name", "john", userId}
	mockError := errors.New("failed to update")
	repositoryMock := mock.EventRepositoryMock{
		CreateKeyFunc: func(ctx context.Context, event *model.EventSnapshot) error {
			if *event == createEvent {
				return mockError
			}
			return nil
		},
	}

	service := NewEventService(&repositoryMock)

	err := service.CreateKey(ctx, &createEvent)

	assertions.ShouldContain(err.Error(), mockError.Error())
}
