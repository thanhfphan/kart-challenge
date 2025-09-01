package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// OutboxEventBuilder helps create outbox events with consistent structure
type OutboxEventBuilder struct {
	eventType     string
	aggregateID   string
	aggregateType string
	eventData     interface{}
	version       int
}

// NewOutboxEventBuilder creates a new builder for outbox events
func NewOutboxEventBuilder() *OutboxEventBuilder {
	return &OutboxEventBuilder{
		aggregateType: AggregateTypeOrder,
		version:       1,
	}
}

// WithEventType sets the event type
func (b *OutboxEventBuilder) WithEventType(eventType string) *OutboxEventBuilder {
	b.eventType = eventType
	return b
}

// WithAggregateID sets the aggregate ID
func (b *OutboxEventBuilder) WithAggregateID(aggregateID string) *OutboxEventBuilder {
	b.aggregateID = aggregateID
	return b
}

// WithAggregateType sets the aggregate type
func (b *OutboxEventBuilder) WithAggregateType(aggregateType string) *OutboxEventBuilder {
	b.aggregateType = aggregateType
	return b
}

// WithEventData sets the event data
func (b *OutboxEventBuilder) WithEventData(eventData interface{}) *OutboxEventBuilder {
	b.eventData = eventData
	return b
}

// WithVersion sets the event version
func (b *OutboxEventBuilder) WithVersion(version int) *OutboxEventBuilder {
	b.version = version
	return b
}

// Build creates the outbox event
func (b *OutboxEventBuilder) Build() (*OutboxEvent, error) {
	eventDataJSON, err := json.Marshal(b.eventData)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	return &OutboxEvent{
		EventType:     b.eventType,
		AggregateID:   b.aggregateID,
		AggregateType: b.aggregateType,
		EventData:     datatypes.JSON(eventDataJSON),
		Status:        OutboxEventStatusPending,
		Version:       b.version,
		CreatedAt:     now,
	}, nil
}

// CreateOrderPlacedEvent creates an outbox event for order placed
func CreateOrderPlacedEvent(order *Order, items []*OrderItem, products []*Product) (*OutboxEvent, error) {
	eventData := NewOrderPlacedEventData(order, items, products)

	return NewOutboxEventBuilder().
		WithEventType(EventTypeOrderPlaced).
		WithAggregateID(order.ID).
		WithAggregateType(AggregateTypeOrder).
		WithEventData(eventData).
		Build()
}

// CreateOrderCompletedEvent creates an outbox event for order completed
func CreateOrderCompletedEvent(order *Order, items []*OrderItem, products []*Product) (*OutboxEvent, error) {
	eventData := NewOrderCompletedEventData(order, items, products)

	return NewOutboxEventBuilder().
		WithEventType(EventTypeOrderCompleted).
		WithAggregateID(order.ID).
		WithAggregateType(AggregateTypeOrder).
		WithEventData(eventData).
		Build()
}

// CreateOrderCancelledEvent creates an outbox event for order cancelled
func CreateOrderCancelledEvent(order *Order, items []*OrderItem, products []*Product, reason string) (*OutboxEvent, error) {
	eventData := NewOrderCancelledEventData(order, items, products, reason)

	return NewOutboxEventBuilder().
		WithEventType(EventTypeOrderCancelled).
		WithAggregateID(order.ID).
		WithAggregateType(AggregateTypeOrder).
		WithEventData(eventData).
		Build()
}

// CreateGenericOrderEvent creates a generic order event with custom data
func CreateGenericOrderEvent(eventType string, order *Order, customData interface{}) (*OutboxEvent, error) {
	return NewOutboxEventBuilder().
		WithEventType(eventType).
		WithAggregateID(order.ID).
		WithAggregateType(AggregateTypeOrder).
		WithEventData(customData).
		Build()
}

// OutboxEventHelpers provides utility functions for working with outbox events
type OutboxEventHelpers struct{}

// NewOutboxEventHelpers creates a new instance of OutboxEventHelpers
func NewOutboxEventHelpers() *OutboxEventHelpers {
	return &OutboxEventHelpers{}
}

// ParseEventData parses the JSON event data into the specified type
func (h *OutboxEventHelpers) ParseEventData(event *OutboxEvent, target interface{}) error {
	return json.Unmarshal(event.EventData, target)
}

// GetOrderPlacedEventData extracts OrderPlacedEventData from an outbox event
func (h *OutboxEventHelpers) GetOrderPlacedEventData(event *OutboxEvent) (*OrderPlacedEventData, error) {
	var eventData OrderPlacedEventData
	err := h.ParseEventData(event, &eventData)
	if err != nil {
		return nil, err
	}
	return &eventData, nil
}

// GetOrderCompletedEventData extracts OrderCompletedEventData from an outbox event
func (h *OutboxEventHelpers) GetOrderCompletedEventData(event *OutboxEvent) (*OrderCompletedEventData, error) {
	var eventData OrderCompletedEventData
	err := h.ParseEventData(event, &eventData)
	if err != nil {
		return nil, err
	}
	return &eventData, nil
}

// GetOrderCancelledEventData extracts OrderCancelledEventData from an outbox event
func (h *OutboxEventHelpers) GetOrderCancelledEventData(event *OutboxEvent) (*OrderCancelledEventData, error) {
	var eventData OrderCancelledEventData
	err := h.ParseEventData(event, &eventData)
	if err != nil {
		return nil, err
	}
	return &eventData, nil
}

// IsOrderEvent checks if the event is related to orders
func (h *OutboxEventHelpers) IsOrderEvent(event *OutboxEvent) bool {
	return event.AggregateType == AggregateTypeOrder
}

// IsEventType checks if the event matches the specified type
func (h *OutboxEventHelpers) IsEventType(event *OutboxEvent, eventType string) bool {
	return event.EventType == eventType
}
