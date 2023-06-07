package contract

import "event-history/pkg/eventinfo/dto"

const EventInfoCreationSuccess = "key created successfully"
const EventInfoUpdateSuccess = "key updated successfully"

type EventFormatter struct {
	EventResponses *dto.EventResponse
}

func (sf *EventFormatter) FormatEventInfoResponse() interface{} {
	return map[string]interface{}{
		"Key":   sf.EventResponses.Key,
		"Value": sf.EventResponses.Value,
	}
}
