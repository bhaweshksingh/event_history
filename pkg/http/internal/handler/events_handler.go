package handler

import (
	"context"
	"event-history/pkg/eventinfo"
	"event-history/pkg/eventinfo/dto"
	"event-history/pkg/eventinfo/model"
	"event-history/pkg/http/contract"
	"event-history/pkg/http/internal/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"go.uber.org/zap"
)

type EventsHandler struct {
	lgr *zap.Logger
	svc eventinfo.Service
}

func NewEventsHandler(lgr *zap.Logger, svc eventinfo.Service) *EventsHandler {
	return &EventsHandler{
		lgr: lgr,
		svc: svc,
	}
}

func (sih *EventsHandler) Create(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	var eventInfo model.EventSnapshot
	err := utils.ParseRequest(req, &eventInfo)
	if err != nil {
		return err
	}

	err = sih.svc.CreateKey(ctx, &eventInfo)
	if err != nil {
		return fmt.Errorf("EventsHandler.CreateKey . error %v", err)
	}

	sih.lgr.Debug("msg", zap.String("eventCode", contract.EventInfoCreationSuccess))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.EventInfoCreationSuccess)
	return nil
}

func (sih *EventsHandler) Update(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	var eventInfo model.EventSnapshot
	err := utils.ParseRequest(req, &eventInfo)
	if err != nil {
		return err
	}

	err = sih.svc.UpdateKey(ctx, &eventInfo)
	if err != nil {
		return fmt.Errorf("EventsHandler.UpdateKey . error %v", err)
	}

	sih.lgr.Debug("msg", zap.String("eventCode", contract.EventInfoUpdateSuccess))
	utils.WriteSuccessResponse(resp, http.StatusCreated, contract.EventInfoUpdateSuccess)
	return nil
}

func (sih *EventsHandler) Delete(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	key, ok := mux.Vars(req)["key"]
	if !ok || len(key) < 1 {
		return fmt.Errorf("URL query Param 'key' is missing")
	}

	userId, ok := mux.Vars(req)["user_id"]
	if !ok || len(userId) < 1 {
		return fmt.Errorf("URL query Param 'userId' is missing")
	}

	err := sih.svc.DeleteKey(ctx, &dto.EventQuery{key, userId})
	if err != nil {
		return fmt.Errorf("error occurred while fetching key details Infos: %v", err)
	}
	utils.WriteSuccessResponse(resp, http.StatusOK, nil)
	return nil
}

func (sih *EventsHandler) Get(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	key, ok := mux.Vars(req)["key"]
	if !ok || len(key) < 1 {
		return fmt.Errorf("URL query Param 'key' is missing")
	}

	userId, ok := mux.Vars(req)["user_id"]
	if !ok || len(userId) < 1 {
		return fmt.Errorf("URL query Param 'userId' is missing")
	}

	eventResponse, err := sih.svc.GetAnswer(ctx, &dto.EventQuery{key, userId})
	if err != nil {
		return fmt.Errorf("error occurred while fetching key details Infos: %v", err)
	}
	sf := &contract.EventFormatter{eventResponse}
	utils.WriteSuccessResponse(resp, http.StatusOK, sf.FormatEventInfoResponse())
	return nil
}

func (sih *EventsHandler) GetHistory(resp http.ResponseWriter, req *http.Request) error {
	ctx := context.Background()
	key, ok := mux.Vars(req)["key"]
	if !ok || len(key) < 1 {
		return fmt.Errorf("URL query Param 'key' is missing")
	}

	userId, ok := mux.Vars(req)["user_id"]
	if !ok || len(userId) < 1 {
		return fmt.Errorf("URL query Param 'userId' is missing")
	}

	historyResponse, err := sih.svc.GetHistory(ctx, &dto.EventQuery{key, userId})
	if err != nil {
		return fmt.Errorf("error occurred while fetching key details Infos: %v", err)
	}
	utils.WriteSuccessResponse(resp, http.StatusOK, historyResponse)
	return nil
}
