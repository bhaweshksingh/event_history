package repository

import (
	"context"
	"event-history/pkg/config"
	"event-history/pkg/eventinfo/dto"
	"event-history/pkg/eventinfo/model"
	"github.com/smartystreets/assertions"
	"gorm.io/gorm"
	"testing"
)

const userId = "test_user"

func getDBConnection() *gorm.DB {
	dbHandler := NewDBHandler(config.NewConfig("").GetDBConfig())
	db, _ := dbHandler.GetDB()
	return db
}

func setUp() (*gorm.DB, context.Context) {
	dbConn := getDBConnection()
	ctx := context.Background()
	dbConn.WithContext(ctx).Where("user_id = ?", userId).Delete(&model.EventSnapshot{})
	dbConn.WithContext(ctx).Where("user_id = ?", userId).Delete(&model.EventHistory{})
	return dbConn, ctx
}
func TestGormEventRepository_GetAnswer(t *testing.T) {
	dbConn, ctx := setUp()
	expectedEvent := model.EventSnapshot{"name", "john", userId}
	dbConn.WithContext(ctx).Create(expectedEvent)
	repository := NewEventRepository(dbConn)

	eventSnapshot, err := repository.GetAnswer(ctx, &dto.EventQuery{"name", userId})

	assertions.So(err, assertions.ShouldBeNil)
	assertions.ShouldEqual(eventSnapshot, expectedEvent)
	checkForEmptyHistory(dbConn, ctx)
}

func TestGormEventRepository_GetAnswer_fails(t *testing.T) {
	dbConn, ctx := setUp()
	repository := NewEventRepository(dbConn)

	eventSnapshot, err := repository.GetAnswer(ctx, &dto.EventQuery{"name", userId})

	assertions.So(eventSnapshot, assertions.ShouldBeNil)
	assertions.ShouldContain(err.Error(), "record not found")
}

func TestGormEventRepository_DeleteAnswer(t *testing.T) {
	dbConn, ctx := setUp()
	event := model.EventSnapshot{"name", "", userId}
	dbConn.WithContext(ctx).Create(event)
	repository := NewEventRepository(dbConn)

	err := repository.DeleteKey(ctx, &dto.EventQuery{"name", userId})

	assertions.So(err, assertions.ShouldBeNil)
	checkForEmptyHistory(dbConn, ctx)
	var historyRecords []*model.EventHistory
	dbConn.WithContext(ctx).Find(&historyRecords)
	assertions.ShouldEqual(len(historyRecords), 1)
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(&event, model.DeleteAction))
}

func TestGormEventRepository_DeleteAnswer_fails(t *testing.T) {
	dbConn, ctx := setUp()
	repository := NewEventRepository(dbConn)

	err := repository.DeleteKey(ctx, &dto.EventQuery{"name", userId})

	assertions.ShouldContain(err.Error(), "record not found")
	checkForEmptyHistory(dbConn, ctx)
	checkForEmptyHistory(dbConn, ctx)
}

func TestGormEventRepository_UpdateKey(t *testing.T) {
	dbConn, ctx := setUp()
	event := &model.EventSnapshot{"name", "john", userId}
	dbConn.WithContext(ctx).Create(event)
	repository := NewEventRepository(dbConn)
	updateEvent := &model.EventSnapshot{"name", "sam", userId}

	err := repository.UpdateKey(ctx, updateEvent)

	assertions.So(err, assertions.ShouldBeNil)
	var historyRecords []*model.EventHistory
	dbConn.WithContext(ctx).Find(&historyRecords)
	assertions.ShouldEqual(len(historyRecords), 1)
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(updateEvent, model.UpdateAction))
}

func TestGormEventRepository_UpdateKey_fails(t *testing.T) {
	dbConn, ctx := setUp()
	repository := NewEventRepository(dbConn)
	updateEvent := &model.EventSnapshot{"name", "sam", userId}

	err := repository.UpdateKey(ctx, updateEvent)

	assertions.ShouldContain(err.Error(), "record not found")
	checkForEmptyHistory(dbConn, ctx)
	checkForEmptySnapshot(dbConn, ctx)
}

func TestGormEventRepository_CreateKey(t *testing.T) {
	dbConn, ctx := setUp()
	repository := NewEventRepository(dbConn)
	createEvent := &model.EventSnapshot{"name", "sam", userId}

	err := repository.CreateKey(ctx, createEvent)

	assertions.So(err, assertions.ShouldBeNil)

	var historyRecords []*model.EventHistory
	dbConn.WithContext(ctx).Find(&historyRecords)
	assertions.ShouldEqual(len(historyRecords), 1)
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(createEvent, model.CreateAction))
}

func TestGormEventRepository_CreateKey_fails(t *testing.T) {
	dbConn, ctx := setUp()
	event := &model.EventSnapshot{"name", "john", userId}
	dbConn.WithContext(ctx).Create(event)
	repository := NewEventRepository(dbConn)
	createEvent := &model.EventSnapshot{"name", "sam", userId}

	err := repository.CreateKey(ctx, createEvent)

	assertions.ShouldContain(err.Error(), "duplicate key value violates unique constraint")
	checkForEmptyHistory(dbConn, ctx)
}

func TestGormEventRepository_GetHistory(t *testing.T) {
	dbConn, ctx := setUp()
	repository := NewEventRepository(dbConn)
	createEvent := &model.EventSnapshot{"name", "john", userId}
	updateEvent := &model.EventSnapshot{"name", "sam", userId}
	deleteEvent := &dto.EventQuery{"name", userId}

	repository.CreateKey(ctx, createEvent)
	repository.UpdateKey(ctx, updateEvent)
	repository.DeleteKey(ctx, deleteEvent)

	historyRecords, err := repository.GetHistory(ctx, &dto.EventQuery{"name", userId})

	assertions.So(err, assertions.ShouldBeNil)
	assertions.ShouldEqual(len(historyRecords), 3)
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(createEvent, model.CreateAction))
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(updateEvent, model.UpdateAction))
	assertions.ShouldEqual(historyRecords[0], model.NewHistoryRecord(&model.EventSnapshot{Key: "name", UserId: userId}, model.DeleteAction))
}

func checkForEmptyHistory(dbConn *gorm.DB, ctx context.Context) {
	var historyRecords []*model.EventHistory
	dbConn.WithContext(ctx).Find(&historyRecords)
	assertions.ShouldEqual(len(historyRecords), 0)
}

func checkForEmptySnapshot(dbConn *gorm.DB, ctx context.Context) {
	var snapshot []*model.EventSnapshot
	dbConn.WithContext(ctx).Find(&snapshot)
	assertions.ShouldEqual(len(snapshot), 0)
}
