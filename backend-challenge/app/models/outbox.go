package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
)

// OutboxEvent represents an event stored in the outbox table for reliable message publishing
type OutboxEvent struct {
	ID            int64          `gorm:"primaryKey;column:id" json:"id"`
	EventType     string         `gorm:"column:event_type;type:varchar(100);not null" json:"event_type"`
	AggregateID   string         `gorm:"column:aggregate_id;type:varchar(36);not null" json:"aggregate_id"`
	AggregateType string         `gorm:"column:aggregate_type;type:varchar(50);not null;default:order" json:"aggregate_type"`
	EventData     datatypes.JSON `gorm:"column:event_data;type:json;not null" json:"event_data"`
	Status        string         `gorm:"column:status;type:varchar(20);not null;default:pending" json:"status"`
	Version       int            `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt     int64          `gorm:"column:created_at;not null" json:"created_at"`
	ProcessedAt   *int64         `gorm:"column:processed_at" json:"processed_at,omitempty"`
}

// OutboxEvent status constants
const (
	OutboxEventStatusPending   = "pending"
	OutboxEventStatusProcessed = "processed"
	OutboxEventStatusFailed    = "failed"
)

// OutboxEvent types constants
const (
	EventTypeOrderPlaced    = "order.placed"
	EventTypeOrderCompleted = "order.completed"
	EventTypeOrderCancelled = "order.cancelled"
)

// Aggregate types constants
const (
	AggregateTypeOrder = "order"
)

func (*OutboxEvent) TableName() string {
	return "outbox_events"
}

func (oe *OutboxEvent) CacheKey() string {
	return fmt.Sprintf("outbox_events:%d", oe.ID)
}

// MarkAsProcessed marks the outbox event as processed
func (oe *OutboxEvent) MarkAsProcessed() {
	oe.Status = OutboxEventStatusProcessed
	now := time.Now().Unix()
	oe.ProcessedAt = &now
}

// MarkAsFailed marks the outbox event as failed
func (oe *OutboxEvent) MarkAsFailed() {
	oe.Status = OutboxEventStatusFailed
}

// IsPending returns true if the event is in pending status
func (oe *OutboxEvent) IsPending() bool {
	return oe.Status == OutboxEventStatusPending
}

// IsProcessed returns true if the event has been processed
func (oe *OutboxEvent) IsProcessed() bool {
	return oe.Status == OutboxEventStatusProcessed
}
