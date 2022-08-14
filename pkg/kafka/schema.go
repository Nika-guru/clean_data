package kafka

type EventData interface {
	any
}

type Schema[T EventData] struct {
	Id        int    `json:"id"`
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
	Data      T      `json:"data"`
}

func (schema *Schema[T]) ContainsEventTypes(eventTypes []string) bool {
	for _, eventType := range eventTypes {
		if eventType == schema.EventType {
			return true
		}
	}
	return false
}
