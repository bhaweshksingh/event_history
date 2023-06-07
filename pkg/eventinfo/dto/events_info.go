package dto

import (
	"event-history/pkg/eventinfo/model"
)

type EventQuery struct {
	Key    string
	UserId string
}

type EventResponse struct {
	Key   string
	Value string
}

// mapping and formatting happens here
func NewEventResponse(eventSnapshot *model.EventSnapshot) *EventResponse {
	return &EventResponse{
		Key:   eventSnapshot.Key,
		Value: eventSnapshot.Value,
	}
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventHistoryResponse struct {
	Data  Data   `json:"data"`
	Event string `json:"event"`
}

func NewEventHistoryResponse(history []model.EventHistory) []EventHistoryResponse {
	var historyResponse []EventHistoryResponse
	for _, event := range history {
		historyResponse = append(historyResponse, EventHistoryResponse{Data{event.Key, event.Value}, event.Action})
	}
	return historyResponse
}
